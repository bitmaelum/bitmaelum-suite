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

	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

// GetAccountTicket retrieves a ticket that can be used for uploading a message
func (api *API) GetAccountTicket(from, to hash.Hash, subscriptionID string) (*ticket.Ticket, error) {
	url := fmt.Sprintf("/account/%s/ticket", from.String())

	return api.retrieveTicket(url, from, to, subscriptionID)
}

// GetTicket retrieves a ticket that can be used for uploading a message
func (api *API) GetTicket(from, to hash.Hash, subscriptionID string) (*ticket.Ticket, error) {
	return api.retrieveTicket("/ticket", from, to, subscriptionID)
}

func (api *API) retrieveTicket(url string, from, to hash.Hash, subscriptionID string) (*ticket.Ticket, error) {
	data, err := json.MarshalIndent(jsonOut{
		"sender":          from.String(),
		"recipient":       to.String(),
		"subscription_id": subscriptionID,
		"preference":      []string{"pow"},
	}, "", "  ")
	if err != nil {
		return nil, err
	}

	body, statusCode, err := api.Post(url, data)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		logrus.Trace(string(body))
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, GetErrorFromResponse(body)
	}

	// Parse body for ticket
	t := &ticket.Ticket{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// ValidateTicket will send back a ticket-result
func (api *API) ValidateTicket(from, to hash.Hash, subscriptionID string, t *ticket.Ticket) (*ticket.Ticket, error) {
	data, err := json.Marshal(jsonOut{
		t.Work.GetName(): t.Work.GetWorkProofOutput(),
	})
	if err != nil {
		return nil, err
	}

	body, statusCode, err := api.Post(fmt.Sprintf("/ticket/%s", t.ID), data)
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		if isErrorResponse(body) {
			return nil, GetErrorFromResponse(body)
		}
		return nil, errNoSuccess
	}

	if isErrorResponse(body) {
		return nil, GetErrorFromResponse(body)
	}

	// Parse body for ticket
	newT, err := ticket.NewFromBytes(body)
	if err != nil {
		return nil, err
	}

	return newT, nil
}
