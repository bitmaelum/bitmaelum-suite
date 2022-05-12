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

package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/store"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// StoreEntry is the unencrypted store entry with child nodes if needed
type StoreEntry struct {
	Path      hash.Hash    `json:"path"`
	Parent    *hash.Hash   `json:"parent"`
	Data      string       `json:"data"`
	Timestamp int64        `json:"timestamp"`
	Children  []StoreEntry `json:"children"`
}

// StoreGetPath will fetch an entry
func (api *API) StoreGetPath(key bmcrypto.KeyPair, addr hash.Hash, path string, recursive bool, since time.Time) (*StoreEntry, error) {
	pathHash := hash.New(addr.String() + path)

	recursiveStr := "0"
	if recursive {
		recursiveStr = "1"
	}
	sinceStr := "0"
	if !since.IsZero() {
		sinceStr = fmt.Sprintf("%d", since.Unix())
	}

	body, statusCode, err := api.Get(fmt.Sprintf("/account/%s/store/%s?recursive=%s&since=%s", addr.String(), pathHash.String(), recursiveStr, sinceStr))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	entry := &store.EntryType{}
	err = json.Unmarshal(body, &entry)
	if err != nil {
		return nil, err
	}

	return decryptStoreEntry(key, entry)
}

func decryptStoreEntry(key bmcrypto.KeyPair, entry *store.EntryType) (*StoreEntry, error) {
	// Validate signature
	if !validateSignature(key.PubKey, entry) {
		fmt.Println("cannot validate signature: ", entry.Path)
		return nil, errNoSuccess
	}

	// Decrypt message
	data, err := bmcrypto.MessageDecrypt(deriveAesKey(key.PrivKey), entry.Data)
	if err != nil {
		fmt.Println("cannot decrypt: ", entry.Path)
		return nil, errNoSuccess
	}

	storeEntry := StoreEntry{
		Path:      entry.Path,
		Parent:    entry.Parent,
		Data:      string(data),
		Timestamp: entry.Timestamp,
		Children:  []StoreEntry{},
	}

	// Iterate all children
	for _, child := range entry.Children {
		tmp, err := decryptStoreEntry(key, &child)
		if err == nil {
			storeEntry.Children = append(storeEntry.Children, *tmp)
		}
	}

	return &storeEntry, nil
}

// StorePutValue will store an value to a path
func (api *API) StorePutValue(key bmcrypto.KeyPair, addr hash.Hash, path string, value string) error {
	// Calc parentPath
	parentPath, _ := filepath.Split(path)
	// correct "/foo/" to "/foo" from "/foo/bar"
	parentPath = strings.TrimRight(parentPath, "/")
	// correct "" to "/" from "/foo"
	if parentPath == "" {
		parentPath = "/"
	}

	pathHash := hash.New(addr.String() + path)
	parentHash := hash.New(addr.String() + parentPath)
	var ph *hash.Hash = &parentHash

	var parent interface{} = parentHash.String()
	if path == "/" {
		parent = nil
		ph = nil
	}

	encValue, err := bmcrypto.MessageEncrypt(deriveAesKey(key.PrivKey), []byte(value))
	if err != nil {
		return errNoSuccess
	}

	sig, err := generateSignature(key.PrivKey, pathHash, ph, encValue)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(jsonOut{
		"path":       pathHash,
		"parent":     parent,
		"value":      encValue,
		"signature":  sig,
		"public_key": key.PubKey.String(),
	}, "", "  ")
	if err != nil {
		return err
	}

	_, statusCode, err := api.Post(fmt.Sprintf("/account/%s/store/%s", addr.String(), pathHash.String()), data)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

func deriveAesKey(privKey bmcrypto.PrivKey) []byte {
	sha := sha256.New()
	sha.Write(privKey.B)
	return sha.Sum(nil)
}

func validateSignature(pubKey bmcrypto.PubKey, entry *store.EntryType) bool {
	return true
}

func generateSignature(privKey bmcrypto.PrivKey, keyHash hash.Hash, parentHash *hash.Hash, value []byte) ([]byte, error) {
	sha := sha256.New()
	sha.Write(keyHash.Byte())
	if parentHash != nil {
		sha.Write(parentHash.Byte())
	}
	sha.Write(value)
	out := sha.Sum(nil)

	return bmcrypto.Sign(privKey, out)
}
