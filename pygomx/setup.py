# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import os
import subprocess

from setuptools import Command, setup
from setuptools.command.build import build
from setuptools.command.bdist_wheel import bdist_wheel


class bdist_wheel_abi3(bdist_wheel):
    def get_tag(self):
        python, abi, plat = super().get_tag()

        if python.startswith("cp") and not (python.endswith("t") or abi.endswith("t")):
            if "android" in plat:
                # cibuildwheel supports android since cp313, so we can't mark it as 310
                return python, "abi3", plat
            # On CPython, our wheels are abi3 and compatible back to 3.10.
            # Free-threaded builds ("t" tag) must keep their original tags (PEP 803).
            # Once PEP 803 is accepted, we may be able to build abi3t wheels.
            return "cp310", "abi3", plat

        return python, abi, plat


class CustomCommand(Command):
    def initialize_options(self) -> None:
        pass

    def finalize_options(self) -> None:
        pass

    def run(self) -> None:
        # configure the go build via env vars
        # keep defaults in sync with build_ffi.py
        go_tags_str = os.getenv("PYGOMX_GO_TAGS")
        if go_tags_str and len(go_tags_str.strip()) > 0:
            go_tags = go_tags_str.split(",")
        else:
            go_tags = []

        match os.getenv("PYGOMX_BUILD_MODE", "static"):
            case "static":
                build_mode_name = "c-archive"
                if os.name == "nt":
                    build_mode_ext = ".lib"
                else:
                    build_mode_ext = ".a"
            case "shared":
                build_mode_name = "c-shared"
                if os.name == "nt":
                    build_mode_ext = ".dll"
                else:
                    build_mode_ext = ".so"
            case _:
                raise ValueError("Invalid PYGOMX_BUILD_MODE.")

        match os.getenv("PYGOMX_OLM_FLAVOR", "goolm"):
            case "none":
                go_tags += ["nocrypto"]
            case "colm":
                go_tags += ["colm"]
            case "goolm":
                go_tags += ["goolm"]
            case "vodozemac":
                go_tags += ["vodozemac"]
                raise ValueError("Vodozemac not supported (yet).")
            case _:
                raise ValueError("Invalid PYGOMX_OLM_FLAVOR.")

        go_call = [
            "go",
            "build",
            f"-buildmode={build_mode_name}",
            "-tags",
            ",".join(go_tags),
            "-o",
            f"../pygomx/libmxclient{build_mode_ext}",
            ".",
        ]
        # print(f"DEBUG: {' '.join(go_call) }")
        subprocess.check_call(go_call, cwd="../libmxclient")

        if os.name == "nt" and os.getenv("PYGOMX_BUILD_MODE", "nope") == "shared":
            from setuptools._distutils.compilers.C.msvc import Compiler

            comp = Compiler()
            comp.initialize()
            subprocess.check_call(
                [
                    comp.lib,
                    "/def:libmxclient.def",
                    "/machine:AMD64",
                    "/out:libmxclient.lib",
                ]
            )


class CustomBuild(build):
    sub_commands = [("build_custom", None)] + build.sub_commands


setup(
    use_calver="%Y.%m.%d.%H.%M",
    setup_requires=["calver"],
    cffi_modules=["build_ffi.py:ffibuilder"],
    cmdclass={
        "build": CustomBuild,
        "build_custom": CustomCommand,
        "bdist_wheel": bdist_wheel_abi3,
    },
)
