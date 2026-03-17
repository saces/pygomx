# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import asyncio
import json
import logging
import threading

from _pygomx import ffi, lib

from .apiv0 import ApiV0Api
from .errors import CheckApiError, CheckApiResult, PygomxAPIError

logger = logging.getLogger(__name__)


class _AsyncClient:
    """
    core binding
    """

    def __init__(self):
        super().__init__()
        self._createMXClient()
        # create a c-handle for self and keep it alive
        self._ffi_selfhandle = ffi.new_handle(self)

        r = lib.apiv0_set_on_event_handler(
            self.client_id, on_event_callback, self._ffi_selfhandle
        )
        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise PygomxAPIError(result)

        r = lib.apiv0_set_on_message_handler(
            self.client_id, on_message_callback, self._ffi_selfhandle
        )
        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise PygomxAPIError(result)

        r = lib.apiv0_set_on_sys_handler(
            self.client_id, on_sys_callback, self._ffi_selfhandle
        )
        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise PygomxAPIError(result)

    def _createMXClient(self):
        r = lib.apiv0_createclient_pass(b".mxpass", b".", b"*", b"*", b"*")

        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise PygomxAPIError(result)

        result_dict = json.loads(result)
        self.client_id = result_dict["id"]
        self.UserID = result_dict["userid"]
        self.DeviceID = result_dict["deviceid"]

    async def _sync(self):
        r = ApiV0Api.startclient(self.client_id)
        CheckApiError(r)

    def _stopsync(self):
        r = ApiV0Api.stopclient(self.client_id)
        CheckApiError(r)

    async def _sendmessage(self, data_dict):
        r = ApiV0Api.sendmessage(self.client_id, data_dict)
        return CheckApiResult(r)

    def leaveroom(self, roomid):
        r = ApiV0Api.leaveroom(self.client_id, roomid)
        CheckApiError(r)

    async def joinedrooms(self):
        r = ApiV0Api.joinedrooms(self.client_id)
        return CheckApiResult(r)

    def _createroom(self, data_dict):
        r = ApiV0Api.createroom(self.client_id, data_dict)
        return CheckApiError(r)

    def process_event(self, evt):
        if hasattr(self, "on_event") and callable(self.on_event):
            loop = asyncio.new_event_loop()
            threading.Thread(
                target=loop.run_forever, name="Async Runner", daemon=True
            ).start()
            asyncio.run_coroutine_threadsafe(self.on_event(evt), loop).result()
        else:
            logger.warn(f"got event but on_event not declared: {evt}")

    def process_message(self, msg):
        if hasattr(self, "on_message") and callable(self.on_message):
            loop = asyncio.new_event_loop()
            threading.Thread(
                target=loop.run_forever, name="Async Runner", daemon=True
            ).start()
            asyncio.run_coroutine_threadsafe(self.on_message(msg), loop).result()
        else:
            logger.warn(f"got message but on_message not declared: {msg}")

    def process_sys(self, ntf):
        if hasattr(self, "on_sys") and callable(self.on_sys):
            loop = asyncio.new_event_loop()
            threading.Thread(
                target=loop.run_forever, name="Async Runner", daemon=True
            ).start()
            asyncio.run_coroutine_threadsafe(self.on_sys(ntf), loop).result()
        else:
            logger.warn(f"got systen notification but on_sys not declared: {ntf}")


@ffi.callback("void(char*, void*)")
def on_event_callback(evt, pobj):
    cli = ffi.from_handle(pobj)
    e = ffi.string(evt).decode("utf-8")
    evt = json.loads(e)
    cli.process_event(evt)


@ffi.callback("void(char*, void*)")
def on_message_callback(msg, pobj):
    cli = ffi.from_handle(pobj)
    m = ffi.string(msg).decode("utf-8")
    msg = json.loads(m)
    cli.process_message(msg)


@ffi.callback("void(char*, void*)")
def on_sys_callback(msg, pobj):
    cli = ffi.from_handle(pobj)
    m = ffi.string(msg).decode("utf-8")
    sys = json.loads(m)
    cli.process_sys(sys)
