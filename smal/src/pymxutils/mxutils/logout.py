# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import getpass
from datetime import datetime

import click
from pygomx.errors import PygomxAPIError

from pygomx import CliV0

from .click import click_catch_exception


@click.command()
@click.option("-u", "--url", "hs_url", metavar="url", help="homeserver url")
@click.option("-t", "--token", "token", metavar="token", help="access token")
@click.option(
    "--json", "show_json", is_flag=True, help="show json as returned from server."
)
@click.argument("devices", metavar="DeviceID", nargs=-1)
@click.option(
    "-l",
    "--logout",
    "logout_type",
    type=click.Choice(["all", "other", "self"]),
    help="logout devices",
)
@click_catch_exception(handle=(PygomxAPIError))
def logout(hs_url, token, devices, logout_type, show_json):
    """List or logout devices.

    \b
    mxlogout [--json]
        list all devices
    mxlogout --logout all
        logout all devices
    mxlogout --logout self
        logout this device
    mxlogout --logout other
        logout all other devices (requires auth)
    mxlogout deviceid [deviceid]...
        logout given devices (requires auth)
    """

    if len(devices) > 0 and logout_type is not None:
        raise ValueError("you can't get both.")

    if hs_url is None and token is None:
        cli = CliV0.from_mxpass(".mxpass", "*", "*", "*")
    else:
        cli = CliV0(hs_url, token)

    match logout_type:
        case "self":
            do_logout(cli, all=False)
            return
        case "all":
            do_logout(cli, all=True)
            return

    if len(devices) > 0:
        device_list = list(devices)
        whoami_dict = cli.Whoami()
        self_user_id = whoami_dict["user_id"]
        do_logout_devices(cli, device_list, self_user_id)
        return

    reqData = {"method": "GET", "path": ["_matrix", "client", "v3", "devices"]}
    raw_device_dict = cli.Generic(reqData)

    if logout_type == "other":
        whoami_dict = cli.Whoami()
        self_device_id = whoami_dict["device_id"]
        self_user_id = whoami_dict["user_id"]
        device_list = []
        for device in raw_device_dict["devices"]:
            if device["device_id"] != self_device_id:
                device_list += [
                    device["device_id"],
                ]
        if len(device_list) > 0:
            do_logout_devices(cli, device_list, self_user_id)
        return

    if show_json:
        click.echo(raw_device_dict)
        return

    max_len = 0
    for device in raw_device_dict["devices"]:
        max_len = max(max_len, len(device["device_id"]))

    for device in raw_device_dict["devices"]:
        date_object = datetime.fromtimestamp(device["last_seen_ts"] / 1000)
        click.echo(
            " ".join(
                [
                    device["device_id"],
                    " " * (max_len - len(device["device_id"])),
                    str(date_object),
                    device["last_seen_ip"],
                    device["display_name"],
                ]
            )
        )
        date_object = datetime.fromtimestamp(device["last_seen_ts"] / 1000)


def do_logout(cli, all):
    reqData = {
        "method": "POST",
        "path": ["_matrix", "client", "v3", "logout"],
    }
    if all:
        reqData["path"] += ["all"]
    res = cli.Generic(reqData)
    click.echo(res)


def do_logout_devices(cli, devices, user_id):
    reqData = {
        "method": "POST",
        "path": ["_matrix", "client", "v3", "delete_devices"],
        "payload": {
            "devices": devices,
            "auth": {
                "type": "m.login.password",
                "password": getpass.getpass(prompt="Password: "),
                "identifier": {
                    "type": "m.id.user",
                    "user": user_id,
                },
            },
        },
    }
    res = cli.Generic(reqData)
    click.echo(res)
