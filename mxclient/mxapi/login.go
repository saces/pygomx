// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxapi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

type LoginInfo struct {
	DiscoverInfo DiscoverInfo `json:"discover_info"`
	Password     string       `json:"password"`
	DeviceID     id.DeviceID  `json:"deviceid"`
	DeviceName   string       `json:"devicename"`
	MXPassFile   string       `json:"mxpassfile"`
}

func Login(li *LoginInfo) (string, error) {

	if li.MXPassFile != "" {
		if _, err := os.Stat(li.MXPassFile); err == nil {
			return "", fmt.Errorf("mxpassfile '%s' already exists", li.MXPassFile)
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("error while checking mxpassfile: %v", err)
		}
	}

	mauclient, err := mautrix.NewClient(li.DiscoverInfo.Homeserver, li.DiscoverInfo.UserID, "")
	if err != nil {
		return "", err
	}

	now := time.Now()
	if li.DeviceID == "" {
		li.DeviceID = id.DeviceID(fmt.Sprintf("libmxclient-%d", now.Unix()))
	}

	if li.DeviceName == "" {
		li.DeviceName = fmt.Sprintf("libmxclient-%s", now.Format(time.RFC3339))
	}

	resp, err := mauclient.Login(context.Background(), &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: li.DiscoverInfo.LoginName,
		},
		Password:                 li.Password,
		DeviceID:                 li.DeviceID,
		InitialDeviceDisplayName: li.DeviceName,
		StoreCredentials:         false,
		StoreHomeserverURL:       false,
		RefreshToken:             false,
	})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("%s | %s | %s | %s\n", li.DiscoverInfo.Homeserver, li.DiscoverInfo.LoginName, li.DiscoverInfo.UserID.Homeserver(), resp.AccessToken)

	if li.MXPassFile != "" {
		err := os.WriteFile(li.MXPassFile, []byte(res), 0600)
		if err != nil {
			return "", fmt.Errorf("unable to write file: %w", err)
		}
		return "", nil
	}

	return res, nil
}
