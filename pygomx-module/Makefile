SOURCE_GO=$(shell find ../libmxclient/ -iname "*.go")

.PHONY:
all: _pygomx.o libmxclient.h
	mkdir -p build/lib
	cp *.so build/lib

_pygomx.o: libmxclient.h build_ffi.py
	python3 build_ffi.py

libmxclient.h: ../libmxclient/go.mod ../libmxclient/go.sum $(SOURCE_GO)
	cd ../libmxclient/ && \
	CGO_ENABLED=1 go build -buildmode=c-archive -o ../pygomx-module/libmxclient.a .

.PHONY:
clean:
	-rm _pygomx.*
	-rm libmxclient.*
	-rm -r dist/
	-rm -r build/
	-rm -r pygomx_module.egg-info

.PHONY:
install: all
	python setup.py bdist_wheel
	python setup.py install
