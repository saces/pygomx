# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from setuptools import setup


setup(
    use_calver="%Y.%m.%d.%H.%M",
    setup_requires=["calver"],
)
