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

	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateApiKey Create a new API key
func (api *API) CreateApiKey(addrHash hash.Hash, key apikey.KeyType) error {
	// Zero is not 1970, but year 1
	var expires int64
	if ! key.Expires.IsZero() {
		expires = key.Expires.Unix()
	}

	data, err := json.MarshalIndent(jsonOut{
		"permissions": key.Permissions,
		"expires":     expires,
		"description": key.Desc,
	}, "", "  ")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/account/%s/apikey", addrHash.String())
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

// DeleteApiKey deletes a new API key
func (api *API) DeleteApiKey(addrHash hash.Hash, ID string) error {
	url := fmt.Sprintf("/account/%s/apikey/%s", addrHash.String(), ID)
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

// ListApiKeys lists alll API keys
func (api *API) ListApiKeys(addrHash hash.Hash) ([]apikey.KeyType, error) {
	url := fmt.Sprintf("/account/%s/apikey", addrHash.String())
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
	keys := &[]apikey.KeyType{}
	err = json.Unmarshal(body, &keys)
	if err != nil {
		return nil, err
	}

	return *keys, nil
}

// GetApiKey gets a single key
func (api *API) GetApiKey(addrHash hash.Hash, ID string) (*apikey.KeyType, error) {
	url := fmt.Sprintf("/account/%s/apikey/%s", addrHash.String(), ID)
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
	key := &apikey.KeyType{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
