#!/usr/bin/python
from cffi import FFI

ffibuilder = FFI()

ffibuilder.set_source(
    module_name="_pygomx",
    source=""" //passed to the real C compiler
        #include "libmxclient.h"
    """,
    libraries=["mxclient"],
)

ffibuilder.cdef(
    csource="""
    extern char* hello(char* p0);
    """
)

if __name__ == "__main__":
    ffibuilder.compile(verbose=True)
