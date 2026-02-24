# set some defaults

# one of static, shared
BUILD_MODE ?= static

# one of none, colm, goolm, vodozemac
OLM_FLAVOR ?= colm

ifndef GO_TAGS
  GO_TAGS =
else
  GO_TAGS := $(GO_TAGS),
endif

ifeq ($(OLM_FLAVOR),none)
    GO_TAGS := $(GO_TAGS)nocrypto
else ifeq ($(OLM_FLAVOR),colm)
    GO_TAGS := $(GO_TAGS)colm
else ifeq ($(OLM_FLAVOR),goolm)
    GO_TAGS := $(GO_TAGS)goolm
else ifeq ($(OLM_FLAVOR),vodozemac)
    GO_TAGS := $(GO_TAGS)vodozemac
else
  $(error unknown OLM_FLAVOR $(OLM_FLAVOR))
endif


# calculate configuration matrix

ifeq ($(BUILD_MODE),static)
  build_mode_name=c-archive
  build_mode_ext=a
else ifeq ($(BUILD_MODE),shared)
  build_mode_name=c-shared
  build_mode_ext=so
else
  $(error unknown BUILD_MODE $(BUILD_MODE))
endif

link_list =

ifeq ($(BUILD_MODE),static)
  ifeq ($(OLM_FLAVOR),colm)
    link_list += olm
  endif
endif


SOURCE_GO=$(shell find ../libmxclient/ -iname "*.go")

ifndef GO_OUTDIR
  GO_OUTDIR = ../pygomx-module
endif

.PHONY:
all: _pygomx.o libmxclient.h
	mkdir -p build/lib
	cp *.so build/lib

_pygomx.o: libmxclient.h build_ffi.py
	python3 build_ffi.py $(link_list)

libmxclient.h: ../libmxclient/go.mod ../libmxclient/go.sum $(SOURCE_GO)
	cd ../libmxclient/ && \
	CGO_ENABLED=1 go build -buildmode=$(build_mode_name) -tags $(GO_TAGS) -o $(GO_OUTDIR)/libmxclient.$(build_mode_ext) .

.PHONY:
clean:
	@rm -f _pygomx.*
	@rm -f libmxclient.*
	@rm -rf dist/
	@rm -rf build/
	@rm -rf pygomx_module.egg-info

.PHONY:
install: all
	python setup.py bdist_wheel
	python setup.py install
