# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import logging

from .app import SMALApp
from .config import SMALConfig

logger = logging.getLogger(__name__)

"""

"""


class SMALBot(SMALApp):
    """ """

    def __init__(self, sigil, config=None):
        super().__init__()
        self._sigil = sigil
        self._config = config if config is not None else SMALConfig.create_new()

    async def sendmessage(self, roomid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.text"

        return await self._sendmessage(data)

    async def sendmessagereply(self, roomid, msgid, mxid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.text"
        data["content"]["m.mentions"] = {}
        data["content"]["m.mentions"]["user_ids"] = [
            mxid,
        ]
        data["content"]["m.relates_to"] = {"m.in_reply_to": {"event_id": msgid}}

        return await self._sendmessage(data)

    async def sendmessagestartthread(self, roomid, msgid, mxid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.text"
        data["content"]["m.mentions"] = {}
        data["content"]["m.mentions"]["user_ids"] = [
            mxid,
        ]
        data["content"]["m.relates_to"] = {"rel_type": "m.thread", "event_id": msgid}
        return await self._sendmessage(data)

    async def sendnoticereply(self, roomid, msgid, mxid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.notice"
        data["content"]["m.relates_to"] = {}
        data["content"]["m.relates_to"]["m.in_reply_to"] = {}
        data["content"]["m.relates_to"]["m.in_reply_to"]["event_id"] = msgid

        return await self._sendmessage(data)

    async def sendnotice(self, roomid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.notice"

        return await self._sendmessage(data)
