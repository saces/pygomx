# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import logging

from _pygomx import lib

from .util import _stringresult, _autostring, _autodict
from .errors import CheckApiError, CheckApiErrorOnly, CheckApiResult

logger = logging.getLogger(__name__)


class ApiV0Api:
    """api_v0 api
    some c-api wrappers & helpers
    inputs: str or bytes
    output: str (utf8)
    """

    @staticmethod
    def discover(mxid):
        return _stringresult(lib.apiv0_discover(_autostring(mxid)))

    @staticmethod
    def login(data):
        return _stringresult(lib.apiv0_login(_autodict(data)))

    @staticmethod
    def joinedrooms(cid):
        return _stringresult(lib.apiv0_joinedrooms(cid))

    @staticmethod
    def sendmessage(cid, data):
        return _stringresult(lib.apiv0_sendmessage(cid, _autodict(data)))

    @staticmethod
    def startclient(cid):
        return _stringresult(lib.apiv0_startclient(cid))

    @staticmethod
    def stopclient(cid):
        return _stringresult(lib.apiv0_stopclient(cid))

    @staticmethod
    def leaveroom(cid, roomid):
        return _stringresult(lib.apiv0_leaveroom(cid, _autostring(roomid)))

    @staticmethod
    def createroom(cid, data):
        return _stringresult(lib.apiv0_createroom(cid, _autodict(data)))

    @staticmethod
    def generic(cid, method, path, data):
        return _stringresult(
            lib.apiv0_genericrequest(
                cid, _autostring(method), _autolist(path), _autodict(data)
            )
        )


class ApiV0:
    """ApiV0"""

    @staticmethod
    def Discover(mxid):
        res = ApiV0Api.discover(mxid)
        return CheckApiResult(res)

    @staticmethod
    def Login(data: dict, mxpassfile: str):
        withpass = False
        if mxpassfile is not None and len(mxpassfile.strip()) > 0:
            data["mxpassfile"] = mxpassfile
            withpass = True
        res = ApiV0Api.login(data)
        if withpass:
            CheckApiError(res)
        else:
            CheckApiErrorOnly(res)
            return res
