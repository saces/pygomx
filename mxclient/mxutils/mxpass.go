// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxutils

import (
	"errors"
	"mxclient/determinant/mxpassfile"
)

func GetMXPassItem(mxpassfile_path string, url string, localpart string, domain string) (*mxpassfile.Entry, error) {
	pf, err := mxpassfile.ReadPassfile(mxpassfile_path)
	if err != nil {
		return nil, err
	}
	e := pf.FindPasswordFill(url, localpart, domain)
	if e == nil {
		return nil, errors.New("no item found.")
	}
	return e, nil
}
