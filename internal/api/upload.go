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
	"io"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
)

// UploadHeader uploads a header
func (api *API) UploadHeader(t ticket.Ticket, header *message.Header) error {
	api.setTicketHeader(t)
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}

	resp, statusCode, err := api.Post("/incoming/header", data)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// UploadCatalog uploads a catalog
func (api *API) UploadCatalog(t ticket.Ticket, encryptedCatalog []byte) error {
	api.setTicketHeader(t)
	_, statusCode, err := api.Post("/incoming/catalog", encryptedCatalog)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// UploadBlock uploads a message block or attachment
func (api *API) UploadBlock(t ticket.Ticket, blockID string, r io.Reader) error {
	api.setTicketHeader(t)
	url := fmt.Sprintf("/incoming/block/%s", blockID)
	_, statusCode, err := api.PostReader(url, r)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// DeleteMessage deletes a message and all content
func (api *API) DeleteMessage(t ticket.Ticket) error {
	api.setTicketHeader(t)

	_, statusCode, err := api.Delete("/incoming")
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// CompleteUpload signals the mailserver that all blocks (and headers) have been uploaded and can start processing
func (api *API) CompleteUpload(t ticket.Ticket) error {
	api.setTicketHeader(t)
	_, statusCode, err := api.Post("/incoming", []byte{})
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}
