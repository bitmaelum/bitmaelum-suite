// Copyright (c) 2022 BitMaelum Authors
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

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/dispatcher"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

var (
	errAuthKeyNotFound = errors.New("auth key not found")
)

type inputAuthKeyType struct {
	Fingerprint string           `json:"fingerprint"`
	PublicKey   *bmcrypto.PubKey `json:"public_key"`
	Signature   string           `json:"signature"`
	Expires     int64            `json:"expires,omitempty"`
	Desc        string           `json:"description,omitempty"`
}

// CreateAuthKey is a handler that will create a new auth key
func CreateAuthKey(w http.ResponseWriter, req *http.Request) {
	var input inputAuthKeyType
	err := httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	//
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	newAuthKey := key.NewAuthKey(*h, input.PublicKey, input.Signature, time.Unix(input.Expires, 0), input.Desc)

	// Store auth key into persistent storage
	repo := container.Instance.GetAuthKeyRepo()
	err = repo.Store(newAuthKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		httputils.ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	_ = dispatcher.DispatchAuthKeyCreate(*h, newAuthKey)

	// Output key
	_ = httputils.JSONOut(w, http.StatusCreated, jsonOut{
		"auth_key": newAuthKey.Fingerprint,
	})
}

// ListAuthKeys returns a list of all keys for the given account
func ListAuthKeys(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	repo := container.Instance.GetAuthKeyRepo()
	keys, err := repo.FetchByHash(h.String())
	if err != nil {
		msg := fmt.Sprintf("error while retrieving authn keys: %s", err)
		httputils.ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	_ = httputils.JSONOut(w, http.StatusOK, keys)
}

// DeleteAuthKey will remove a key
func DeleteAuthKey(w http.ResponseWriter, req *http.Request) {
	k, err := hasAuthKeyAccess(w, req)
	if err != nil {
		return
	}

	repo := container.Instance.GetAuthKeyRepo()
	_ = repo.Remove(*k)

	_ = dispatcher.DispatchAuthKeyDelete(k.AddressHash, *k)

	// All is well
	_ = httputils.JSONOut(w, http.StatusNoContent, "")
}

// GetAuthKeyDetails will get a key
func GetAuthKeyDetails(w http.ResponseWriter, req *http.Request) {
	k, err := hasAuthKeyAccess(w, req)
	if err != nil {
		return
	}

	// Output key
	_ = httputils.JSONOut(w, http.StatusOK, k)
}

func hasAuthKeyAccess(w http.ResponseWriter, req *http.Request) (*key.AuthKeyType, error) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errAccountNotFound.Error())
		return nil, errAccountNotFound
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.Instance.GetAuthKeyRepo()
	k, err := repo.Fetch(keyID)
	if err != nil || k.AddressHash.String() != h.String() {
		httputils.ErrorOut(w, http.StatusNotFound, errAuthKeyNotFound.Error())
		return nil, errAuthKeyNotFound
	}

	return k, nil
}
