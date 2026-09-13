// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only

//go:build !windows

package mxpassfile

import (
	"errors"
	"os"
)

// ReadPassfile reads the file at path and parses it into a Passfile.
func ReadPassfile(path string) (*Passfile, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	permissions := fileInfo.Mode().Perm()
	if permissions != 0o600 {
		return nil, errors.New("To wide permissions, ignore file")
	}
	return readPassfile(path)
}
