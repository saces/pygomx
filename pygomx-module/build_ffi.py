#!/usr/bin/python
import os
import sys
from cffi import FFI

lib_list = [
    "mxclient",
]

# keep defaults in sync with setup.py
if (
    os.getenv("PYGOMX_BUILD_MODE", "static") == "static"
    and os.getenv("PYGOMX_OLM_FLAVOR", colm) == "colm"
):
    lib_list += ["olm"]

print(f"liblist: {lib_list}")

ffibuilder = FFI()

ffibuilder.set_source(
    module_name="_pygomx",
    source=""" //passed to the real C compiler
        #include "libmxclient.h"
    """,
    libraries=lib_list,
    library_dirs=["."],
    include_dirs=["."],
)

ffibuilder.cdef(
    csource="""
    typedef void (*on_event_handler_ptr) (char*, void*);
    typedef void (*on_message_handler_ptr) (char*, void*);
    typedef void (*on_sys_handler_ptr) (char*, void*);
    extern void  FreeCString(char* s);
    extern char* cli_discoverhs(char* mxid);
    extern char* cli_mkmxtoken(char* mxid, char* pw);
    extern char* cli_whoami(char* hs, char* accessToken);
    extern char* cli_accountinfo(char* hs, char* accessToken);
    extern char* cli_clearaccount(char* hs, char* accessToken);
    extern char* cli_serverinfo(char* url);
    extern int   apiv0_initialize();
    extern int   apiv0_deinitialize();
    extern char* apiv0_discover(char* mxid);
    extern char* apiv0_login(char* data);
    extern char* apiv0_createclient(char* storage_path, char* hs, char* mxid, char* accessToken);
    extern char* apiv0_createclient_pass(char* mxpassfile, char* storage_path, char* hs, char* localpart, char* domain);   
    extern char* apiv0_set_on_event_handler(int cid, on_event_handler_ptr ptr, void* pobj);
    extern char* apiv0_set_on_message_handler(int cid, on_message_handler_ptr ptr, void* pobj);
    extern char* apiv0_set_on_sys_handler(int cid, on_sys_handler_ptr ptr, void* pobj);
    extern char* apiv0_startclient(int cid);
    extern char* apiv0_stopclient(int cid);
    extern char* apiv0_sendmessage(int cid, char* data);
    extern char* apiv0_leaveroom(int cid, char* roomid);
    extern char* apiv0_joinedrooms(int cid);
    extern char* apiv0_createroom(int cid, char* data);
    extern int   apiv0_removeclient(int cid);
    extern char* apiv0_listclients();
    extern char* apiv0_getoptions(int cid);
    extern int   apiv0_setoptions(int cid, char* opts);
    """
)

if __name__ == "__main__":
    ffibuilder.compile(verbose=True)
