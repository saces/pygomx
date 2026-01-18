import sys
from _pygomx import lib, ffi


def mktoken():
    if len(sys.argv) != 3:
        print("usage: ", sys.argv[0], "  matrixid password")
        return 1

    mxid = sys.argv[1].encode(encoding="utf-8")
    pw = sys.argv[2].encode(encoding="utf-8")

    r = lib.cli_mkmxtoken(mxid, pw)
    result = ffi.string(r)
    print(result)
