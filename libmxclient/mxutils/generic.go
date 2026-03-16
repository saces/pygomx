// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxutils

import (
	"context"
	"encoding/json"

	"maunium.net/go/mautrix"
)

type req_data struct {
	Method  string              `json:"method"`
	Path    mautrix.BaseURLPath `json:"path"`
	Payload any                 `json:"payload"`
}

func GenericRequest(hs string, accessToken string, reqData string) (string, error) {
	var rd req_data
	err := json.Unmarshal([]byte(reqData), &rd)
	if err != nil {
		return "", err
	}
	mauclient, err := mautrix.NewClient(hs, "", accessToken)
	if err != nil {
		return "", err
	}

	urlPath := mauclient.BuildURLWithFullQuery(mautrix.BaseURLPath(rd.Path), nil)
	resp, err := mauclient.MakeRequest(context.Background(), rd.Method, urlPath, rd.Payload, nil)
	return string(resp), err
}
