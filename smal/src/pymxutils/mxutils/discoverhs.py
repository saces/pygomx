import sys
from _pygomx import lib, ffi


def discoverhs():
    if len(sys.argv) != 2:
        print("usage: ", sys.argv[0], "  matrixid|domainname")
        return 1

    mxid = sys.argv[1].encode(encoding="utf-8")

    print("try to discover from: ", mxid)

    r = lib.cli_discoverhs(mxid)
    result = ffi.string(r)
    print(result)
