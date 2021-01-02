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
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

type inputInviteType struct {
	AddrHash string `json:"address"`
	Days     int    `json:"days"`
}

type jsonOut map[string]interface{}

// NewInvite handler will generate a new invite token for a given address
func NewInvite(w http.ResponseWriter, req *http.Request) {
	k := handler.GetAPIKey(req)
	if !k.HasPermission(internal.PermGenerateInvites, nil) {
		httputils.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input inputInviteType
	err := httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	addrHash, err := hash.NewFromHash(input.AddrHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect address")
		return
	}

	validUntil := internal.TimeNow().Add(time.Duration(input.Days) * 24 * time.Hour)
	token, err := signature.NewInviteToken(*addrHash, config.Routing.RoutingID, validUntil, config.Routing.PrivateKey)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "cannot generate invite token")
		return
	}

	_ = httputils.JSONOut(w, http.StatusCreated, jsonOut{
		"hash":   addrHash.String(),
		"token":  token.String(),
		"expiry": validUntil.Unix(),
	})
}
