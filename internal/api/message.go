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
	"fmt"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

// Message is a standard structure that returns a message header + catalog
type Message struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"h"`
	Catalog []byte         `json:"c"`
}

// GetMessage retrieves a message header + catalog from a message box
func (api *API) GetMessage(addr hash.Hash, box, messageID string) (*Message, error) {
	in := &Message{}

	url := fmt.Sprintf("/account/%s/box/%s/message/%s", addr.String(), box, messageID)
	resp, statusCode, err := api.GetJSON(url, in)
	if err != nil {
		logrus.Trace(err)
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, getErrorFromResponse(resp)
	}

	return in, nil
}

// GetMessageBlock retrieves a message block
func (api *API) GetMessageBlock(addr hash.Hash, box, messageID, blockID string) ([]byte, error) {
	body, statusCode, err := api.Get(fmt.Sprintf("/account/%s/box/%s/message/%s/block/%s", addr.String(), box, messageID, blockID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return body, nil
}

// GetMessageAttachment retrieves a message attachment reader
func (api *API) GetMessageAttachment(addr hash.Hash, box, messageID, attachmentID string) (io.ReadCloser, error) {
	r, statusCode, err := api.GetReader(fmt.Sprintf("/account/%s/box/%s/message/%s/attachment/%s", addr.String(), box, messageID, attachmentID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return r, nil
}
