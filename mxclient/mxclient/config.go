// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxclient

type ClientCreateConfig struct {
	MXPassfilePath string `json:"passfile_path"`
	DBType         string `json:"db_type"`
	DBName         string `json:"db_name"`
}
