# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import logging
import asyncio

from pygomx.client import _AsyncClient

logger = logging.getLogger(__name__)

"""

"""


class SMALApp(_AsyncClient):
    """

    implement 'async def self.on_startup()'
        async_client is logged in & ready.
        time to setup extra things & hooks not covered by this class
        sync_loop will not start til we return

    implement 'async def self.on_startup_run()'
        async_client is logged in & ready.
        this will not wait for return
        do your even long running startup code here

    """

    def __init__(self):
        super().__init__()

    def run(self, sync=True):
        asyncio.run(self.main_loop(sync))

    async def main_loop(self, sync):
        if hasattr(self, "on_startup") and callable(self.on_startup):
            await self.on_startup()

        if hasattr(self, "on_startup_run") and callable(self.on_startup_run):
            await asyncio.ensure_future(self.on_startup_run())

        if sync:
            await self._sync()

    def stop(self):
        self._stopsync()
