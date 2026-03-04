import os
import subprocess
from contextlib import suppress
from setuptools import Command, setup
from setuptools.command.build import build


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
                build_mode_ext = ".a"
            case "shared":
                build_mode_name = "c-shared"
                build_mode_ext = ".so"
            case _:
                raise ValueError("Invalid PYGOMX_BUILD_MODE.")

        match os.getenv("PYGOMX_OLM_FLAVOR", "colm"):
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
            f"-buildmode=${build_mode_name}",
            "-tags",
            ",".join(go_tags),
            "-o",
            f"../pygomx-module/libmxclient${build_mode_ext}",
            ".",
        ]
        subprocess.call(go_call, cwd="../libmxclient")


class CustomBuild(build):
    sub_commands = [("build_custom", None)] + build.sub_commands


setup(
    cffi_modules=["build_ffi.py:ffibuilder"],
    cmdclass={"build": CustomBuild, "build_custom": CustomCommand},
)
