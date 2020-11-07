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

package handlers

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"golang.org/x/sync/errgroup"
)

// ComposeMessage composes a new message from the given account Info to the "to" with given subject, blocks and attachments
func ComposeMessage(addressing message.Addressing, subject string, b, a []string) error {
	envelope, err := message.Compose(addressing, subject, b, a)
	if err != nil {
		return err
	}

	// Setup API connection to the server
	client, err := api.NewAuthenticated(addressing.Sender.Address, addressing.Sender.PrivKey, addressing.Sender.Host)
	if err != nil {
		return err
	}

	// and finally send
	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}
	return nil
}

func uploadToServer(addressing message.Addressing, envelope *message.Envelope) error {
	client, err := api.NewAuthenticated(addressing.Sender.Address, addressing.Sender.PrivKey, addressing.Sender.Host)
	if err != nil {
		return err
	}

	// Get upload ticket
	t, err := client.GetTicket(envelope.Header.From.Addr, envelope.Header.To.Addr, "")
	if err != nil {
		return errors.New("cannot get ticket from server: " + err.Error())
	}
	if !t.Valid {
		return errors.New("invalid ticket returned by server")
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
