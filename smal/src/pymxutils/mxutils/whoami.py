# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import click
from pygomx import CliV0
from pygomx.errors import PygomxAPIError

from .click import click_catch_exception


@click.command()
@click.option("-u", "--url", "hs_url", metavar="url", help="homeserver url")
@click.option("-t", "--token", "token", metavar="token", help="access token")
@click_catch_exception(handle=(PygomxAPIError))
def whoami(hs_url, token):
    """this token belongs to?"""

    if hs_url is None and token is None:
        cli = CliV0.from_mxpass(".mxpass", "*", "*", "*")
    else:
        cli = CliV0(hs_url, token)

    click.echo(cli.Whoami())
