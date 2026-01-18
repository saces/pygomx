import sys
from _pygomx import lib, ffi


def accountinfo():
    if len(sys.argv) != 3:
        print("usage: ", sys.argv[0], "  url acesstoken")
        return 1

    url = sys.argv[1].encode(encoding="utf-8")
    tk = sys.argv[1].encode(encoding="utf-8")

    r = lib.cli_accountinfo(url, tk)
    result = ffi.string(r)
    print(result)
