# -*- coding: utf-8 -*-
import logging
from _pygomx import lib, ffi
import json
from .errors import APIError

logger = logging.getLogger(__name__)


def checkApiError(cstr):
    result = ffi.string(cstr).decode("utf-8")
    lib.FreeCString(cstr)

    if result.startswith("ERR:"):
        raise APIError(result)

    if result == "SUCCESS.":
        return

    logger.debug(result)

    result_dict = json.loads(result)
    return result_dict


class _MXClient:
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
            raise APIError(result)

        r = lib.apiv0_set_on_message_handler(
            self.client_id, on_message_callback, self._ffi_selfhandle
        )
        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise APIError(result)

        r = lib.apiv0_set_on_sys_handler(
            self.client_id, on_sys_callback, self._ffi_selfhandle
        )
        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise APIError(result)

    def _createMXClient(self):
        r = lib.apiv0_createclient_pass(b".mxpass", b".", b"*", b"*", b"*")

        result = ffi.string(r)
        lib.FreeCString(r)
        if result.startswith(b"ERR:"):
            raise APIError(result)

        result_dict = json.loads(result)
        self.client_id = result_dict["id"]
        self.UserID = result_dict["userid"]
        self.DeviceID = result_dict["deviceid"]

    def _sync(self):
        r = lib.apiv0_startclient(self.client_id)
        checkApiError(r)

    def _stopsync(self):
        r = lib.apiv0_stopclient(self.client_id)
        checkApiError(r)

    def _sendmessage(self, data_dict):
        data = json.dumps(data_dict).encode(encoding="utf-8")
        r = lib.apiv0_sendmessage(self.client_id, data)
        result = checkApiError(r)
        return result

    def leaveroom(self, roomid):
        r = lib.apiv0_leaveroom(self.client_id, roomid.encode(encoding="utf-8"))
        checkApiError(r)

    def joinedrooms(self):
        r = lib.apiv0_joinedrooms(self.client_id)
        return checkApiError(r)

    def process_event(self, evt):
        if hasattr(self, "on_event") and callable(self.on_event):
            self.on_event(evt)
        else:
            logger.warn(f"got event but on_event not declared: {evt}")

    def process_message(self, msg):
        if hasattr(self, "on_message") and callable(self.on_message):
            self.on_message(msg)
        else:
            logger.warn(f"got message but on_message not declared: {msg}")

    def process_sys(self, ntf):
        if hasattr(self, "on_sys") and callable(self.on_sys):
            self.on_sys(ntf)
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
