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

package api

import (
	"bytes"
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

// RemoveMessage deletes a message on the server
func (api *API) RemoveMessage(addr hash.Hash, messageID string) error {
	url := fmt.Sprintf("/account/%s/message/%s", addr.String(), messageID)
	resp, statusCode, err := api.Delete(url)
	if err != nil {
		logrus.Trace(err)
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// RemoveMessageFromBox deletes a message on the server
func (api *API) RemoveMessageFromBox(addr hash.Hash, messageID string, boxID int) error {
	type inputRemoveMessage struct {
		From int `json:"from"`
	}

	input := &inputRemoveMessage{
		From: boxID,
	}

	url := fmt.Sprintf("/account/%s/message/%s/delete", addr.String(), messageID)
	resp, statusCode, err := api.PostJSON(url, input)
	if err != nil {
		logrus.Trace(err)
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// CopyMessage copies a message to a mailbox
func (api *API) CopyMessage(addr hash.Hash, messageID string, to int) error {
	type inputCopyMessage struct {
		To int `json:"to"`
	}

	input := &inputCopyMessage{
		To: to,
	}

	url := fmt.Sprintf("/account/%s/message/%s/copy", addr.String(), messageID)
	resp, statusCode, err := api.PostJSON(url, input)
	if err != nil {
		logrus.Trace(err)
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// MoveMessage moves a message from one mailbox to another
func (api *API) MoveMessage(addr hash.Hash, messageID string, from, to int) error {
	type inputMoveMessage struct {
		From int `json:"from"`
		To   int `json:"to"`
	}

	input := &inputMoveMessage{
		From: from,
		To:   to,
	}

	url := fmt.Sprintf("/account/%s/message/%s/move", addr.String(), messageID)
	resp, statusCode, err := api.PostJSON(url, input)
	if err != nil {
		logrus.Trace(err)
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// GetMessage retrieves a message header + catalog from a message
func (api *API) GetMessage(addr hash.Hash, messageID string) (*Message, error) {
	in := &Message{}

	url := fmt.Sprintf("/account/%s/message/%s", addr.String(), messageID)
	resp, statusCode, err := api.GetJSON(url, in)
	if err != nil {
		logrus.Trace(err)
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, GetErrorFromResponse(resp)
	}

	return in, nil
}

// GetMessageBlock retrieves a message block
func (api *API) GetMessageBlock(addr hash.Hash, messageID, blockID string) (io.ReadCloser, error) {
	r, statusCode, err := api.GetReader(fmt.Sprintf("/account/%s/message/%s/block/%s", addr.String(), messageID, blockID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return r, nil
}

// GetMessageAttachment retrieves a message attachment reader
func (api *API) GetMessageAttachment(addr hash.Hash, messageID, attachmentID string) (io.ReadCloser, error) {
	r, statusCode, err := api.GetReader(fmt.Sprintf("/account/%s/message/%s/attachment/%s", addr.String(), messageID, attachmentID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return r, nil
}

// GenerateAPIBlockReader returns a reader function that will create a reader from the given message and block ID.
func (api *API) GenerateAPIBlockReader(addr hash.Hash) func(messageID, blockID string) io.Reader {
	return func(messageID, blockID string) io.Reader {
		r, err := api.GetMessageBlock(addr, messageID, blockID)
		if err != nil {
			return bytes.NewReader([]byte{})
		}

		return r
	}
}

// GenerateAPIAttachmentReader returns a reader function that will create a reader from the given message and attachment ID.
func (api *API) GenerateAPIAttachmentReader(addr hash.Hash) func(messageID, attachmentID string) io.Reader {
	return func(messageID, attachmentID string) io.Reader {
		r, err := api.GetMessageAttachment(addr, messageID, attachmentID)
		if err != nil {
			return bytes.NewReader([]byte{})
		}

		return r
	}
}
