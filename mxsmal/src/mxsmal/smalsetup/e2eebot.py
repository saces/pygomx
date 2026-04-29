# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import logging
from mxsmal.bot import SMALBot
from pygomx.apiv0 import ApiV0Api

logger = logging.getLogger(__name__)


# TODO this should be an App, not a Bot
class E2eeBot(SMALBot):

    def __init__(self):
        super().__init__("¿")

    async def on_startup(self):
        ApiV0Api.self_sign(self.client_id)

    async def on_sys(self, ntf):
        print("Got a system notification: ", ntf)

    async def on_event(self, evt):
        print("Got an event: ", evt)

    async def on_message(self, msg):
        print("Got a mesage: ", msg)
