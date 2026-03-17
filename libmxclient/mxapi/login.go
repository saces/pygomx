// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxapi

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

type login_data struct {
	Homeserver      string `json:"homeserver"`
	Mxid            string `json:"mxid"`
	Loginname       string `json:"loginname"`
	Password        string `json:"password"`
	DeviceID        string `json:"deviceid"`
	DeviceName      string `json:"devicename"`
	MXPassFile      string `json:"mxpassfile"`
	MakeMasterKey   bool   `json:"make_master_key"`
	MakeRecoveryKey bool   `json:"make_recovery_key"`
}

func Login(data string) (string, error) {
	var ld login_data
	err := json.Unmarshal([]byte(data), &ld)
	if err != nil {
		return "", err
	}

	if ld.MXPassFile != "" {
		if _, err := os.Stat(ld.MXPassFile); err == nil {
			return "", fmt.Errorf("mxpassfile '%s' already exists", ld.MXPassFile)
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("error while checking mxpassfile: %v", err)
		}
	}

	mauclient, err := mautrix.NewClient(ld.Homeserver, id.UserID(ld.Mxid), "")
	if err != nil {
		return "", err
	}

	now := time.Now()
	if ld.DeviceID == "" {
		ld.DeviceID = fmt.Sprintf("libmxclient-%d", now.Unix())
	}

	if ld.DeviceName == "" {
		ld.DeviceName = fmt.Sprintf("libmxclient-%s", now.Format(time.RFC3339))
	}

	resp, err := mauclient.Login(context.Background(), &mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: "m.id.user",
			User: ld.Loginname,
		},
		Password:                 ld.Password,
		DeviceID:                 id.DeviceID(ld.DeviceID),
		InitialDeviceDisplayName: ld.DeviceName,
		StoreCredentials:         false,
		StoreHomeserverURL:       false,
		RefreshToken:             false,
	})
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("%s | %s | %s | %s\n", ld.Homeserver, ld.Loginname, id.UserID(ld.Mxid).Homeserver(), resp.AccessToken)

	if ld.MakeMasterKey {
		masterkey := make([]byte, 32)
		rand.Read(masterkey)
		res = fmt.Sprintf("%smaster | | | %x\n", res, masterkey)
	}

	if ld.MakeRecoveryKey {
		recoverykey := make([]byte, 32)
		rand.Read(recoverykey)
		res = fmt.Sprintf("%srecovery | | | %x\n", res, recoverykey)
	}

	if ld.MXPassFile != "" {
		err := os.WriteFile(ld.MXPassFile, []byte(res), 0600)
		if err != nil {
			return "", fmt.Errorf("unable to write file: %w", err)
		}
		return "SUCCESS.", nil
	}

	return res, nil
}
