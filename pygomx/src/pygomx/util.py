# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from _pygomx import ffi, lib
import json


def _stringresult(cstr):
    result = ffi.string(cstr)
    lib.FreeCString(cstr)
    return result.decode("utf-8")


def _autostring(xstr):
    match xstr:
        case bytes():
            return xstr
        case str():
            return xstr.encode(encoding="utf-8")
        case _:
            raise TypeError("only str or bytes allowed")


def _autodict(xdict):
    match xdict:
        case bytes():
            return xdict
        case str():
            return xdict.encode(encoding="utf-8")
        case dict():
            return json.dumps(xdict).encode(encoding="utf-8")
        case _:
            raise TypeError("only str or bytes or dict allowed")
