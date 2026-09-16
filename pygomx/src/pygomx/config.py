# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
from dataclasses import dataclass
from enum import StrEnum

class DBType(StrEnum):
    NONE   = "none"
    SQLITE = "sqlite"
    POSTGRES = "postgres"

@dataclass
class PygomxCreateConfig:
    """config object for passing parameters to create"""

    passfile_path: str = ".mxpass"

    db_type: DBType = DBType.SQLITE
    db_name: str = "pygomx.db"
