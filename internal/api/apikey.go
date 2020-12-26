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

	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateAPIKey Create a new API key
func (api *API) CreateAPIKey(addrHash hash.Hash, key key.APIKeyType) error {
	// Zero is not 1970, but year 1
	var expires int64
	if !key.Expires.IsZero() {
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
		return GetErrorFromResponse(body)
	}

	return nil
}

// DeleteAPIKey deletes a new API key
func (api *API) DeleteAPIKey(addrHash hash.Hash, ID string) error {
	url := fmt.Sprintf("/account/%s/apikey/%s", addrHash.String(), ID)
	body, statusCode, err := api.Delete(url)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	if isErrorResponse(body) {
		return GetErrorFromResponse(body)
	}

	return nil
}

// ListAPIKeys lists alll API keys
func (api *API) ListAPIKeys(addrHash hash.Hash) ([]key.APIKeyType, error) {
	url := fmt.Sprintf("/account/%s/apikey", addrHash.String())
	body, statusCode, err := api.Get(url)
	if err != nil {
		return []key.APIKeyType{}, err
	}

	if statusCode < 200 || statusCode > 299 {
		return []key.APIKeyType{}, errNoSuccess
	}

	if isErrorResponse(body) {
		return []key.APIKeyType{}, GetErrorFromResponse(body)
	}

	// Parse body for keys
	keys := &[]key.APIKeyType{}
	err = json.Unmarshal(body, &keys)
	if err != nil || keys == nil {
		return []key.APIKeyType{}, err
	}

	return *keys, nil
}

// GetAPIKey gets a single key
func (api *API) GetAPIKey(addrHash hash.Hash, ID string) (*key.APIKeyType, error) {
	url := fmt.Sprintf("/account/%s/apikey/%s", addrHash.String(), ID)
	body, statusCode, err := api.Get(url)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, GetErrorFromResponse(body)
	}

	// Parse body for key
	k := &key.APIKeyType{}
	err = json.Unmarshal(body, &k)
	if err != nil {
		return nil, err
	}

	return k, nil
}
