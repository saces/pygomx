// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxapi

import (
	"context"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

type DiscoverInfo struct {
	Homeserver string    `json:"homeserver"`
	UserID     id.UserID `json:"user_id"`
	LoginName  string    `json:"login_name"`
}

// Discover tries to guess the homeserver from the given matrix id
func Discover(userID id.UserID) (di *DiscoverInfo, err error) {
	var tmp_di DiscoverInfo
	tmp_di.LoginName, tmp_di.Homeserver, err = userID.ParseAndValidateRelaxed()
	if err != nil {
		return
	}

	wk, err := mautrix.DiscoverClientAPI(context.Background(), tmp_di.Homeserver)
	if err != nil {
		return
	}
	if wk != nil {
		tmp_di.Homeserver = wk.Homeserver.BaseURL
	}
	tmp_di.UserID = userID
	di = &tmp_di
	return
}
