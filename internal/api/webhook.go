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

	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateWebhook Create a new API key
func (api *API) CreateWebhook(wh webhook.Type) (*webhook.Type, error) {
	// ID is set by the server
	wh.ID = ""
	data, err := json.Marshal(wh)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/account/%s/webhook", wh.Account.String())
	body, statusCode, err := api.Post(url, data)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, GetErrorFromResponse(body)
	}

	fmt.Println(body)

	return &wh, nil
}

// DeleteWebhook deletes a webhook
func (api *API) DeleteWebhook(addrHash hash.Hash, ID string) error {
	url := fmt.Sprintf("/account/%s/webhook/%s", addrHash.String(), ID)
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

// ListWebhooks lists alll webhooks
func (api *API) ListWebhooks(addrHash hash.Hash) ([]webhook.Type, error) {
	url := fmt.Sprintf("/account/%s/webhook", addrHash.String())
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

	// Parse body for webhooks
	webhooks := []webhook.Type{}
	err = json.Unmarshal(body, &webhooks)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

// UpdateWebhook will update a webhook
func (api *API) UpdateWebhook(addrHash hash.Hash, ID string, wh webhook.Type) error {
	if wh.ID != ID {
		return errNoSuccess
	}

	data, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/account/%s/webhook/%s", addrHash.String(), ID)
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

// GetWebhook gets a single webhook
func (api *API) GetWebhook(addrHash hash.Hash, ID string) (*webhook.Type, error) {
	url := fmt.Sprintf("/account/%s/webhook/%s", addrHash.String(), ID)
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

	// Parse body for webhook
	wh := &webhook.Type{}
	err = json.Unmarshal(body, &wh)
	if err != nil {
		return nil, err
	}

	return wh, nil
}

// EnableWebhook will enable a webhook
func (api *API) EnableWebhook(addrHash hash.Hash, ID string) error {
	url := fmt.Sprintf("/account/%s/webhook/%s/enable", addrHash.String(), ID)
	body, statusCode, err := api.Post(url, []byte{})
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

// DisableWebhook will disable a webhook
func (api *API) DisableWebhook(addrHash hash.Hash, ID string) error {
	url := fmt.Sprintf("/account/%s/webhook/%s/disable", addrHash.String(), ID)
	body, statusCode, err := api.Post(url, []byte{})
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
