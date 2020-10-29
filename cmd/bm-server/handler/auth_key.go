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

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
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
	err := DecodeBody(w, req.Body, &input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	//
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	newAuthKey := key.NewAuthKey(*h, input.PublicKey, input.Signature, time.Unix(input.Expires, 0), input.Desc)

	// Store Auth key into persistent storage
	repo := container.GetAuthKeyRepo()
	err = repo.Store(newAuthKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"auth_key": newAuthKey.Fingerprint,
	})
}

// ListAuthKeys returns a list of all keys for the given account
func ListAuthKeys(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	// Store Auth key into persistent storage
	repo := container.GetAuthKeyRepo()
	keys, err := repo.FetchByHash(h.String())
	if err != nil {
		msg := fmt.Sprintf("error while retrieving keys: %s", err)
		ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(keys)
}

// DeleteAuthKey will remove a key
func DeleteAuthKey(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.GetAuthKeyRepo()
	k, err := repo.Fetch(keyID)
	if err != nil || k.AddressHash.String() != h.String() {
		// Only allow deleting of keys that we own as account
		ErrorOut(w, http.StatusNotFound, "key not found")
		return
	}

	_ = repo.Remove(*k)

	// All is well
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetAuthKeyDetails will get a key
func GetAuthKeyDetails(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.GetAuthKeyRepo()
	k, err := repo.Fetch(keyID)
	if err != nil || k.AddressHash.String() != h.String() {
		ErrorOut(w, http.StatusNotFound, "key not found")
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(k)
}
