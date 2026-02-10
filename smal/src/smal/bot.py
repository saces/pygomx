# -*- coding: utf-8 -*-
import logging
import sys

from .app import SMALApp

logger = logging.getLogger(__name__)

"""

"""


class SMALBot(SMALApp):
    """ """

    def __init__(self, sigil):
        super().__init__()
        self._sigil = sigil

    def run(self):
        self._sync()

    def stop(self):
        self._stopsync()

    def sendmessage(self, roomid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.text"

        self._sendmessage(data)

    def sendmessagereply(self, roomid, msgid, mxid, text):
        data = {}
        data["roomid"] = roomid
        data["content"] = {}
        data["content"]["body"] = text
        data["content"]["msgtype"] = "m.text"
        data["content"]["m.mentions"] = {}
        data["content"]["m.mentions"]["user_ids"] = [
            mxid,
        ]
        data["content"]["m.relates_to"] = {}
        data["content"]["m.relates_to"]["m.in_reply_to"] = {}
        data["content"]["m.relates_to"]["m.in_reply_to"]["event_id"] = msgid

        self._sendmessage(data)
