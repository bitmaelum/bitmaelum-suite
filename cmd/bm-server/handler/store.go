// Copyright (c) 2021 BitMaelum Authors
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
	"crypto/sha256"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/store"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

var (
	errPathNotFound = errors.New("store: path not found")
)

// StoreGet will retrieve a path or collection
func StoreGet(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	pathHash, err := hash.NewFromHash(mux.Vars(req)["path"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	recursive, since := parseQueryString(req)
	getPath(w, *haddr, *pathHash, recursive, since)
}

func parseQueryString(req *http.Request) (recursive bool, since time.Time) {
	query := req.URL.Query()

	if query.Get("recursive") == "1" {
		recursive = true
	}

	since = time.Time{}
	if query.Get("since") != "" {
		ts, err := strconv.Atoi(query.Get("since"))
		if err == nil {
			since = time.Unix(int64(ts), 0)
		}
	}
	return
}

// UpdateType is a request for a store entry
type UpdateType struct {
	Path      hash.Hash       `json:"path"`
	Parent    *hash.Hash      `json:"parent"`
	Value     []byte          `json:"value"`
	Signature []byte          `json:"signature"`
	PubKey    bmcrypto.PubKey `json:"public_key"`
}

// StoreUpdate will update a path or collection
func StoreUpdate(w http.ResponseWriter, req *http.Request) {
	updateRequest := &UpdateType{}
	err := json.NewDecoder(req.Body).Decode(updateRequest)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "Malformed JSON: "+err.Error())
		return
	}

	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	pathHash, err := hash.NewFromHash(mux.Vars(req)["path"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	// Check if signature matches
	if !checkSignature(updateRequest.PubKey, updateRequest.Path, updateRequest.Parent, updateRequest.Value, updateRequest.Signature) {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	storePath(w, *haddr, *pathHash, updateRequest.Parent, updateRequest.Value, updateRequest.Signature)
}

func checkSignature(pubKey bmcrypto.PubKey, pathHash hash.Hash, parentHash *hash.Hash, value []byte, signature []byte) bool {
	sha := sha256.New()
	sha.Write(pathHash.Byte())
	if parentHash != nil {
		sha.Write(parentHash.Byte())
	}
	sha.Write(value)
	out := sha.Sum(nil)

	ok, err := bmcrypto.Verify(pubKey, out, signature)
	if err != nil {
		return false
	}

	return ok
}

// StoreDelete will remove a path or collection
func StoreDelete(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	pathHash, err := hash.NewFromHash(mux.Vars(req)["path"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	deletePath(w, *haddr, *pathHash)
}

func storePath(w http.ResponseWriter, addrHash, pathHash hash.Hash, parentHash *hash.Hash, value, signature []byte) {
	err := openDb(w, addrHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}
	defer closeDb(addrHash)

	// Add entry
	entry := &store.EntryType{
		Path:      pathHash,
		Parent:    parentHash,
		Data:      value,
		Signature: signature,
	}

	storesvc := container.Instance.GetStoreRepo()
	err = storesvc.SetEntry(addrHash, *entry)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, errPathNotFound.Error())
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, nil)
}

func deletePath(w http.ResponseWriter, addrHash hash.Hash, pathHash hash.Hash) {
	err := openDb(w, addrHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}
	defer closeDb(addrHash)

	// Check if path exists in database
	storesvc := container.Instance.GetStoreRepo()
	if !storesvc.HasEntry(addrHash, pathHash) {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	err = storesvc.RemoveEntry(addrHash, pathHash, false)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, errPathNotFound.Error())
		return
	}

	_ = httputils.JSONOut(w, http.StatusNoContent, nil)
}

func getPath(w http.ResponseWriter, addrHash hash.Hash, pathHash hash.Hash, recursive bool, since time.Time) {
	err := openDb(w, addrHash)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}
	defer closeDb(addrHash)

	// Check if path exists in database
	storesvc := container.Instance.GetStoreRepo()
	if !storesvc.HasEntry(addrHash, pathHash) {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	entry, err := storesvc.GetEntry(addrHash, pathHash, recursive, since)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, entry)
}

func openDb(w http.ResponseWriter, addrHash hash.Hash) error {
	// Open DB
	storesvc := container.Instance.GetStoreRepo()
	if err := storesvc.OpenDb(addrHash); err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, errPathNotFound.Error())
		return errors.New("cannot open db")
	}

	return nil
}

func closeDb(addrhash hash.Hash) {
	storesvc := container.Instance.GetStoreRepo()
	_ = storesvc.CloseDb(addrhash)
}
