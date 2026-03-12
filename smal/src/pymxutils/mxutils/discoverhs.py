# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import sys
from _pygomx import lib, ffi
import click
import json


@click.command()
@click.option(
    "--json", "show_json", is_flag=True, help="show json as returned from server."
)
@click.argument("domain", metavar="string")
def discoverhs(domain, show_json):
    """Attempts to discover the homeserver from the given string"""
    mxid = domain.encode(encoding="utf-8")

    r = lib.cli_discoverhs(mxid)
    result = ffi.string(r).decode("utf-8")
    lib.FreeCString(r)
    if result.startswith("ERR:"):
        print(result)
        sys.exit(1)
    if show_json:
        print(result)
    else:
        result_dict = json.loads(result)
        print(result_dict["m.homeserver"]["base_url"])
