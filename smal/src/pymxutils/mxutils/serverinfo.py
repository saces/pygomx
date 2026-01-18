import sys
from _pygomx import lib, ffi


def serverinfo():
    if len(sys.argv) != 2:
        print("usage: ", sys.argv[0], "  url|domainname")
        return 1

    mxdomain = sys.argv[1].encode(encoding="utf-8")

    r = lib.cli_serverinfo(mxdomain)
    result = ffi.string(r)
    print(result)
