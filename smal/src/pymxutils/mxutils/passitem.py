# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import json

import click
from pygomx.errors import PygomxAPIError

from pygomx import CliV0

from .click import click_catch_exception


@click.command()
@click.option(
    "-s", "--secret", "show_secret", is_flag=True, help="print only the secret"
)
@click.option("-u", "--url", "hs_url", metavar="url", help="url selector")
@click.option(
    "-l", "--localpart", "localpart", metavar="localpart", help="localpart selector"
)
@click.option("-d", "--domain", "domain", metavar="domain", help="domain selector")
@click.argument("mxpassfile", metavar="mxpassfilepath", required=False)
@click_catch_exception(handle=(PygomxAPIError))
def passitem(mxpassfile, show_secret, hs_url, localpart, domain):
    """utility to get items from mxpasss files"""

    # defaults
    if mxpassfile is None:
        mxpassfile = ".mxpass"
    if hs_url is None:
        hs_url = "*"
    if localpart is None:
        localpart = "*"
    if domain is None:
        domain = "*"

    result_dict = CliV0.MXPassItem(mxpassfile, hs_url, localpart, domain)

    if show_secret:
        click.echo(result_dict["Token"])
    else:
        result_dict["Token"] = "***"
        click.echo(json.dumps(result_dict))
