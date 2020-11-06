package messages

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"golang.org/x/sync/errgroup"
)

// Send will send a message inside envelope via routing, using addr/key
func Send(client api.API, envelope *message.Envelope) error {

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
