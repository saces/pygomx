# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import logging

from .app import SMALApp

logger = logging.getLogger(__name__)

"""

"""


class SMALBot(SMALApp):
    """ """

    def __init__(self, config=None):
        super().__init__(config)

    async def sendmessage(self, roomid, text):
        content = {
            "body": text,
            "msgtype": "m.text"
        }
        return await self.room_send_message(roomid, 'm.room.message', content)

    async def sendmessagereply(self, roomid, msgid, mxid, text):
        content = {
            "body":  text,
            "msgtype": "m.text",
            "m.mentions": {
                "user_ids":
                    - mxid
            },
            "m.relates_to":  {"m.in_reply_to": {"event_id": msgid}}
        }
        return await self.room_send_message(roomid, 'm.room.message', content)

    async def sendmessagestartthread(self, roomid, msgid, mxid, text):
        content = {
            "body": text,
            "msgtype":  "m.text",
            "m.mentions": {
                "user_ids":
                  - mxid
            },
            "m.relates_to":  {"rel_type": "m.thread", "event_id": msgid}
        }
        return await self.room_send_message(roomid, 'm.room.message', content)

    async def sendnoticereply(self, roomid, msgid, mxid, text):
        content = {
            "body": text,
            "msgtype": "m.notice",
            "m.relates_to": {
                "m.in_reply_to": {
                    "event_id":  msgid}}}
        return await self.room_send_message(roomid, 'm.room.message', content)

    async def sendnotice(self, roomid, text):
        content = {
            "body": text,
            "msgtype": "m.notice"
        }
        return await self.room_send_message(roomid, 'm.room.message', content)
