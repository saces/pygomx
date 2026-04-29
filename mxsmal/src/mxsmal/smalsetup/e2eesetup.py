# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from datetime import datetime
import getpass
import os
import json
from functools import partial, wraps

import click
from pygomx.errors import PygomxAPIError

from pygomx.apiv0 import ApiV0
from pygomx.cliv0 import CliV0

from .e2eebot import E2eeBot


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
@catch_exception(handle=(PygomxAPIError))
def e2eesetup(mxpassfile):
    """Utility for smalbot e2ee setup"""

    cli = CliV0.from_mxpass(mxpassfile, "*", "*", "*")

    # we need to know who we are and the device id needs to be known too
    whoami = cli.Whoami()

    # check for other devices
    devices = cli.Generic("GET", ["_matrix", "client", "v3", "devices"])

    if len(devices["devices"]) > 1:
        click.echo(f"Other devices found for '{whoami['user_id']}'")
        click.echo(f"This device: {whoami['device_id']}")
        click.echo("Other devices:")
        other_devices_list = []
        for device in devices["devices"]:
            # print(device)
            if device["device_id"] != whoami["device_id"]:
                other_devices_list.append(device["device_id"])
                # devices never logged in don't have a 'last_seen_ts'
                last_seen = ""
                if device["last_seen_ts"]:
                    last_seen = datetime.fromtimestamp(device["last_seen_ts"] / 1000)
                click.echo(
                    f"    {device['device_id']} ({device['display_name']}) - {last_seen} from {device['last_seen_ip']}"
                )
        click.echo()
        if click.confirm("Do you want to log them out?"):

            cli.Generic(
                "POST",
                ["_matrix", "client", "v3", "delete_devices"],
                {
                    "devices": other_devices_list,
                    "auth": {
                        "type": "m.login.password",
                        "password": getpass.getpass(prompt="Password: "),
                        "identifier": {
                            "type": "m.id.user",
                            "user": whoami["user_id"],
                        },
                    },
                },
            )

            click.echo("Well done!")
        # TODO maybe check for master key and reset it too.

    # TODO the config thingy comes here, for now a lot of stuff is hard coded to make it work

    # Start the bot to setup e2ee, the real bot will not be startd

    e2eeBot = E2eeBot()

    e2eeBot.run(sync=False)

    print("Huhu Bämm!")
