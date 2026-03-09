# -*- coding: utf-8 -*-


class APIError(Exception):
    """Exception raised for api usage errors.

    Attributes:
        message -- explanation of the error
    """

    def __init__(self, message):
        self.message = message[4:]
        super().__init__(self.message)
