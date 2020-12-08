// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package mgmt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/dispatcher"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type inputAPIKeyType struct {
	Permissions []string `json:"permissions"`
	Expires     int64    `json:"expires"`
	AddrHash    string   `json:"hash,omitempty"`
	Desc        string   `json:"description,omitempty"`
}

// NewAPIKey is a handler that will create a new API key (non-admin keys only)
func NewAPIKey(w http.ResponseWriter, req *http.Request) {
	var input inputAPIKeyType
	err := httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	// Make sure we can only set api permissions for the account we have permission for.
	var h *hash.Hash
	if input.AddrHash != "" {
		tmp := hash.New(input.AddrHash)
		h = &tmp
	}
	if h == nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect hash")
		return
	}

	k := handler.GetAPIKey(req)
	if !k.HasPermission(internal.PermAPIKeys, h) {
		httputils.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err = internal.CheckManagementPermissions(input.Permissions)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect permissions")
		return
	}

	newAPIKey := key.NewAPIAccountKey(*h, input.Permissions, time.Unix(input.Expires, 0), input.Desc)

	// Store API key into persistent storage
	repo := container.Instance.GetAPIKeyRepo()
	err = repo.Store(newAPIKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		httputils.ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	_ = dispatcher.DispatchApiKeyCreate(*h, newAPIKey)

	// Output key
	_ = httputils.JSONOut(w, http.StatusCreated, jsonOut{
		"api_key": newAPIKey.ID,
	})
}
