# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import click
from pygomx.errors import PygomxAPIError

from pygomx import CliV0

from .click import click_catch_exception


@click.command()
@click.option(
    "--json", "show_json", is_flag=True, help="show json as returned from server."
)
@click.argument("domain", metavar="string")
@click_catch_exception(handle=(PygomxAPIError))
def discoverhs(domain, show_json):
    """Attempts to discover the homeserver from the given string"""
    result = CliV0.Discover(domain)

    if show_json:
        click.echo(result)
    else:
        click.echo(result["m.homeserver"]["base_url"])
