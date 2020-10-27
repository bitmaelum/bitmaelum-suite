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

package api

import (
	"encoding/json"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal/authkey"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateAuthKey will create a new authorized key on the server
func (api *API) CreateAuthKey(addrHash hash.Hash, key authkey.KeyType) error {
	// Zero is not 1970, but year 1
	var expires int64
	if !key.Expires.IsZero() {
		expires = key.Expires.Unix()
	}

	data, err := json.MarshalIndent(jsonOut{
		"fingerprint": key.Fingerprint,
		"public_key":  key.PublicKey,
		"signature":   key.Signature,
		"expires":     expires,
		"description": key.Description,
	}, "", "  ")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/account/%s/authkey", addrHash.String())
	body, statusCode, err := api.Post(url, data)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	if isErrorResponse(body) {
		return getErrorFromResponse(body)
	}

	return nil
}

// DeleteAuthKey deletes a new auth key
func (api *API) DeleteAuthKey(addrHash hash.Hash, fingerprint string) error {
	url := fmt.Sprintf("/account/%s/authkey/%s", addrHash.String(), fingerprint)
	body, statusCode, err := api.Delete(url)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	if isErrorResponse(body) {
		return getErrorFromResponse(body)
	}

	return nil
}

// ListAuthKeys lists all auth keys
func (api *API) ListAuthKeys(addrHash hash.Hash) ([]authkey.KeyType, error) {
	url := fmt.Sprintf("/account/%s/authkey", addrHash.String())
	body, statusCode, err := api.Get(url)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, getErrorFromResponse(body)
	}

	// Parse body for keys
	keys := &[]authkey.KeyType{}
	err = json.Unmarshal(body, &keys)
	if err != nil {
		return nil, err
	}

	return *keys, nil
}

// GetAuthKey gets a single key
func (api *API) GetAuthKey(addrHash hash.Hash, fingerprint string) (*authkey.KeyType, error) {
	url := fmt.Sprintf("/account/%s/authkey/%s", addrHash.String(), fingerprint)
	body, statusCode, err := api.Get(url)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, getErrorFromResponse(body)
	}

	// Parse body for key
	key := &authkey.KeyType{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
