# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from _pygomx import ffi, lib
import json
from .errors import APIError


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


def CheckApiError(rstr):
    if rstr.startswith("ERR:"):
        raise APIError(rstr)

    if rstr == "SUCCESS.":
        return None

    raise ValueError(f"unexpected result: {rstr[:20]}")


def CheckApiResult(rstr):
    if rstr.startswith("ERR:"):
        raise APIError(rstr)

    if rstr == "SUCCESS.":
        return None

    result_dict = json.loads(rstr)
    return result_dict


class CliV0Api:
    """cli_v0 api
    c-api wrappers
    inputs: str or bytes
    output: str (utf8)
    """

    @staticmethod
    def whoami(hs_url, token):
        return _stringresult(lib.cliv0_whoami(_autostring(hs_url), _autostring(token)))

    @staticmethod
    def mxpassitem(mxpassfile, hs_url, localpart, domain):
        return _stringresult(
            lib.cliv0_mxpassitem(
                _autostring(mxpassfile),
                _autostring(hs_url),
                _autostring(localpart),
                _autostring(domain),
            )
        )

    @staticmethod
    def generic(hs_url, token, data):
        return _stringresult(
            lib.cliv0_genericrequest(
                _autostring(hs_url),
                _autostring(token),
                _autodict(data),
            )
        )


class CliV0:
    """cli_v0 api class
    high level & helper functions
    """

    def __init__(self, hs_url, token):
        self.hs_url = hs_url
        self.token = token

    @classmethod
    def from_mxpass(cls, mxpassfile, hs_url, localpart, domain):
        res = CliV0Api.mxpassitem(mxpassfile, hs_url, localpart, domain)
        result_dict = CheckApiResult(res)
        return cls(result_dict["Matrixhost"], result_dict["Token"])

    def Whoami(self):
        res = CliV0Api.whoami(self.hs_url, self.token)
        return CheckApiResult(res)

    def Generic(self, data, omitt_token=False):
        if omitt_token:
            token = ""
        else:
            token = self.token
        res = CliV0Api.generic(self.hs_url, token, data)
        return CheckApiResult(res)
