# libmxclient api doc

C API

simple highlevel matrix api intended for bindings to script languages

`everything is a string`™

this should make bindings easier.

parameters are strings, either a simple name or json strings

return values are either a strings starting with `ERR:` or a json result

some funktions might return `SUCCESS.` as equivalent for `ret 0`

returned strings must be free'ed y the caller

it's the callers task to keep passed callback pointers alive

## versions

`v0` unstable, devolopment version

all other version are stable, may get additions and security fixes, no other changes

