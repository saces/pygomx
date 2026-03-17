# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from _pygomx import lib
from .errors import CheckApiResult

from .util import _stringresult, _autostring, _autodict


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

    @staticmethod
    def discover(domain):
        return _stringresult(
            lib.cliv0_discoverhs(
                _autostring(domain),
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

    @staticmethod
    def Discover(domain):
        res = CliV0Api.discover(domain)
        return CheckApiResult(res)

    @staticmethod
    def MXPassItem(mxpassfile, hs_url, localpart, domain):
        res = CliV0Api.mxpassitem(mxpassfile, hs_url, localpart, domain)
        return CheckApiResult(res)

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
