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

package messages

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"golang.org/x/sync/errgroup"
)

var errInvalidTicket = errors.New("invalid ticket returned by server")

// Send will send a message inside envelope via routing, using addr/key
func Send(client api.API, envelope *message.Envelope) error {

	// Get upload ticket
	t, err := client.GetAccountTicket(envelope.Header.From.Addr, envelope.Header.To.Addr, "")
	if err != nil {
		return err
	}
	if !t.Valid {
		return errInvalidTicket
	}

	// parallelize uploads
	g := new(errgroup.Group)
	g.Go(func() error {
		return client.UploadHeader(*t, envelope.Header)
	})
	g.Go(func() error {
		return client.UploadCatalog(*t, envelope.EncryptedCatalog)
	})
	for id, r := range envelope.BlockReaders {
		// Store locally, otherwise the anonymous go function doesn't know which "id" and "r"
		id := id
		r := r
		g.Go(func() error {
			return client.UploadBlock(*t, id, *r)
		})
	}
	for id, r := range envelope.AttachmentReaders {
		// Store locally, otherwise the anonymous go function doesn't know which "id" and "r"
		id := id
		r := r
		g.Go(func() error {
			return client.UploadBlock(*t, id, *r)
		})
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		_ = client.DeleteMessage(*t)
		return err
	}

	// All done, mark upload as completed
	return client.CompleteUpload(*t)
}
