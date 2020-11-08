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

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

type inputAPIKeyType struct {
	Permissions []string `json:"permissions"`
	Expires     int64    `json:"expires,omitempty"`
	Desc        string   `json:"description,omitempty"`
}

// CreateAPIKey is a handler that will create a new API key (non-admin keys only)
func CreateAPIKey(w http.ResponseWriter, req *http.Request) {
	var input inputAPIKeyType
	err := DecodeBody(w, req.Body, &input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect body")
		return
	}

	//
	err = internal.CheckAccountPermissions(input.Permissions)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect permissions")
		return
	}

	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	newAPIKey := key.NewAPIAccountKey(*h, input.Permissions, time.Unix(input.Expires, 0), input.Desc)

	// Store API key into persistent storage
	repo := container.Instance.GetAPIKeyRepo()
	err = repo.Store(newAPIKey)
	if err != nil {
		msg := fmt.Sprintf("error while storing key: %s", err)
		ErrorOut(w, http.StatusInternalServerError, msg)
		return
	}

	// Output key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(jsonOut{
		"api_key": newAPIKey.ID,
	})
}

// ListAPIKeys returns a list of all keys for the given account
func ListAPIKeys(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	// Store API key into persistent storage
	repo := container.Instance.GetAPIKeyRepo()
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

// DeleteAPIKey will remove a key
func DeleteAPIKey(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.Instance.GetAPIKeyRepo()
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

// GetAPIKeyDetails will get a key
func GetAPIKeyDetails(w http.ResponseWriter, req *http.Request) {
	h, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	keyID := mux.Vars(req)["key"]

	// Fetch key
	repo := container.Instance.GetAPIKeyRepo()
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
