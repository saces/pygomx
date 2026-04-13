# Copyright (C) 2026 saces@c-base.org
# SPDX-License-Identifier: AGPL-3.0-only
import os
import pathlib
import sys
from dataclasses import asdict, dataclass
from enum import Enum

import dacite


class ConfigError(Exception):
    pass


class ConfigFileType(Enum):
    YAML = "yaml"
    TOML = "toml"
    JSON = "json"


@dataclass
class ConfigFile:

    @classmethod
    def create_new(cls):
        return cls()

    @classmethod
    def load_file(cls, file_path=None, file_type=None, cast=None):
        """
        automagic config file loader

        Args:
            file_path: fd, filename or "-"
            type: type of config file, required for fd and "-"
            cast: type hint decite/decoder

        Returns: config object

        """
        if file_path is None:
            return cls()

        if file_type is None:
            # try to guess the type from name extension
            ext = pathlib.Path(file_path).suffix
            if ext in [".yaml", ".yml"]:
                file_type = ConfigFileType.YAML
            elif ext in [".toml", ".tml"]:
                file_type = ConfigFileType.TOML
            elif ext == ".json":
                file_type = ConfigFileType.JSON
            else:
                raise ConfigError("Unable to guess the file type from name.")

        # try to load the data
        data = None
        if file_path == "-":
            file_path = sys.stdin
        if file_type == ConfigFileType.TOML:
            import toml

            data = toml.load(file_path)
        elif file_type == ConfigFileType.YAML:
            import yaml

            data = yaml.safe_load(file_path)
        elif file_type == ConfigFileType.JSON:
            import json

            with open(file_path) as f:
                data = json.load(f)
        else:
            # after adding new file types you may need to update the logic
            raise RuntimeError()

        return dacite.from_dict(
            data_class=cls,
            data=data,
            # tell dacite to cast some types while deserializing
            config=dacite.Config(cast=cast) if cast else None,
        )

    @classmethod
    def load_file_auto(cls, file_path: str, file_type=None, cast=None):
        if os.path.exists(file_path):
            return cls.load_file(file_path, file_type, cast)
        else:
            return cls()

    def save_file(self, file_path: str, file_type=None):

        if file_path == "-" and file_type is None:
            raise ConfigError("file_type required if printing to stdout")

        if file_type is None:
            # try to guess the type from name extension
            ext = pathlib.Path(file_path).suffix
            if ext in [".yaml", ".yml"]:
                file_type = ConfigFileType.YAML
            elif ext in [".toml", ".tml"]:
                file_type = ConfigFileType.TOML
            elif ext == ".json":
                file_type = ConfigFileType.JSON
            else:
                raise ConfigError("Unable to guess the file type from name.")

        def my_factory(data):
            def convert_value(obj):
                if isinstance(obj, Enum):
                    return obj.value
                return obj

            return dict((k, convert_value(v)) for k, v in data)

        d = asdict(self, dict_factory=my_factory)

        if file_path == "-":
            file_path = sys.stdout
            if file_type == ConfigFileType.TOML:
                import toml

                toml.dump(d, file_path)
            elif file_type == ConfigFileType.YAML:
                import yaml

                yaml.dump(d, file_path)
            elif file_type == ConfigFileType.JSON:
                import json

                json.dump(d, file_path)
            else:
                # after adding new file types you may need to update the logic
                raise RuntimeError()
        else:
            if file_type == ConfigFileType.TOML:
                import toml

                with open(file_path, "w") as file:
                    toml.dump(d, file)
            elif file_type == ConfigFileType.YAML:
                import yaml

                with open(file_path, "w") as file:
                    yaml.dump(d, file)
            elif file_type == ConfigFileType.JSON:
                import json

                with open(file_path, "w") as file:
                    json.dump(d, file)


@dataclass
class SMALConfig(ConfigFile):
    pass
