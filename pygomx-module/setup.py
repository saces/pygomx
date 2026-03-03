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
        go_call = [
            "go",
            "build",
            "-buildmode=c-archive",
            "-o",
            "../pygomx-module/libmxclient.a",
            ".",
        ]
        subprocess.call(go_call, cwd="../libmxclient")


class CustomBuild(build):
    sub_commands = [("build_custom", None)] + build.sub_commands


setup(
    cffi_modules=["build_ffi.py:ffibuilder"],
    cmdclass={"build": CustomBuild, "build_custom": CustomCommand},
)
