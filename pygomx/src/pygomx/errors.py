# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import json


class PygomxAPIError(Exception):
    """Exception raised for api usage errors.

    Attributes:
        message -- explanation of the error
    """

    def __init__(self, message):
        self.message = message[4:]
        super().__init__(self.message)


def CheckApiErrorOnly(rstr):
    if rstr.startswith("ERR:"):
        raise PygomxAPIError(rstr)

    return rstr


def CheckApiError(rstr):
    if rstr.startswith("ERR:"):
        raise PygomxAPIError(rstr)

    if rstr == "SUCCESS.":
        return None

    raise ValueError(f"unexpected result: {rstr[:60]}")


def CheckApiResult(rstr):
    if rstr.startswith("ERR:"):
        raise PygomxAPIError(rstr)

    if rstr == "SUCCESS.":
        return None

    result_dict = json.loads(rstr)
    return result_dict
