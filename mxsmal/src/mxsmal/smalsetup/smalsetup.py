# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from datetime import datetime
import getpass
import os
import time
from functools import partial, wraps

import click
from pygomx.errors import PygomxAPIError

from pygomx.apiv0 import ApiV0


def catch_exception(func=None, *, handle):
    if not func:
        return partial(catch_exception, handle=handle)

    @wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except handle as e:
            raise click.ClickException(e)

    return wrapper


@click.command()
@click.option(
    "--mxpass",
    "mxpassfile",
    metavar="filepath",
    default=".mxpass",
    help="mxpass file name",
)
@click.argument("mxid", metavar="MatrixID")
@catch_exception(handle=(PygomxAPIError))
def smalsetup(mxid, mxpassfile):
    """Utility for creating smalbot mxpass files"""

    now = int(time.time())

    if len(mxpassfile.strip()) > 0:
        if os.path.exists(mxpassfile):
            raise click.ClickException(f"file {mxpassfile} exists.")

    discover_info = ApiV0.Discover(mxid)

    login_info = {
        "discover_info": discover_info,
        "deviceid": f"smalbot-{now}",
        "devicename": f"smalbot-{datetime.fromtimestamp(now)}",
        "mxpassfile": mxpassfile,
        "password": getpass.getpass(prompt="Password: "),
    }

    ApiV0.Login(login_info)

    click.echo("login created. start your bot now or run e2eesetup.")
