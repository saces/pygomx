# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import datetime
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

    create_mxpass = len(mxpassfile.strip()) > 0

    if create_mxpass:
        if os.path.exists(mxpassfile):
            raise click.ClickException(f"file {mxpassfile} exists.")

    result_dict = ApiV0.Discover(mxid)

    result_dict["password"] = getpass.getpass(prompt="Password: ")
    result_dict["make_master_key"] = True
    result_dict["make_recovery_key"] = True

    now = int(time.time())
    result_dict["deviceid"] = f"smalbot-{now}"
    result_dict["devicename"] = f"smalbot-{datetime.datetime.fromtimestamp(now)}"

    ApiV0.Login(result_dict, ".mxpass")

    click.echo("login created. start your bot now.")
