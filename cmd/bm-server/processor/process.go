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

package processor

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/dispatcher"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

var errValidating = errors.New("error while validating")

// ProcessMessage will process a message found in the processing queue.
//   * If it's a local address, it will be moved to the local mailbox
//   * If it's a remote address, it will be send to the remote mail server
//   * If things fail, it will be moved to the retry queue, where it will be moved to processed queue later
func ProcessMessage(msgID string) {
	logrus.Debugf("processing message %s", msgID)

	// Set the message in the scoreboard, so we know this message is being processed.
	AddToScoreboard(message.SectionProcessing, msgID)
	defer func() {
		RemoveFromScoreboard(message.SectionProcessing, msgID)
	}()

	// Check header and get recipient
	header, err := message.GetMessageHeader(message.SectionProcessing, msgID)
	if err != nil {
		// cannot read header.. Let's move to retry queue
		logrus.Warnf("cannot find or read header for message %s. Retrying.", msgID)
		MoveToRetryQueue(msgID)
		return
	}

	rs := container.Instance.GetResolveService()
	addrInfo, err := rs.ResolveAddress(header.To.Addr)
	if err != nil {
		logrus.Trace(err)
		logrus.Warnf("cannot resolve address %s for message %s. Retrying.", header.To.Addr, msgID)
		MoveToRetryQueue(msgID)
		return
	}

	// Local addresses don't need to be send. They are treated locally
	ar := container.Instance.GetAccountRepo()
	if ar.Exists(header.To.Addr) {
		// Do stuff locally
		logrus.Debugf("Message %s can be transferred locally to %s", msgID, addrInfo.Hash)

		err := deliverLocal(addrInfo, msgID, header)
		if err != nil {
			logrus.Warnf("cannot deliver message %s locally to %s. Retrying.", msgID, header.To.Addr)
			MoveToRetryQueue(msgID)
		}

		return
	}

	// Deliver remote
	err = deliverRemote(addrInfo, msgID, header)
	if err != nil {
		logrus.Warnf("cannot deliver message %s remotely to %s. Retrying.", msgID, header.To.Addr)
		MoveToRetryQueue(msgID)
	}
}

// deliverLocal moves a message to a local mailbox.
func deliverLocal(addrInfo *resolver.AddressInfo, msgID string, header *message.Header) error {
	// Check the serverSignature
	if !message.VerifyServerHeader(*header) {
		logrus.Errorf("message %s destined for %s has failed the server signature check. Seems that this message did not originate from the original mail server. Removing the message.", msgID, header.To.Addr)

		err := message.RemoveMessage(message.SectionProcessing, msgID)
		if err != nil {
			// @TODO: we should notify somebody?
			logrus.Warnf("cannot remove message %s from the process queue.", msgID)
		}

		return nil
	}

	// Check the clientSignature
	if !message.VerifyClientHeader(*header) {
		logrus.Errorf("message %s destined for %s has failed the client signature check. Seems that this message may have been spoofed. Removing the message.", msgID, header.To.Addr)

		err := message.RemoveMessage(message.SectionProcessing, msgID)
		if err != nil {
			// @TODO: we should notify somebody?
			logrus.Warnf("cannot remove message %s from the process queue.", msgID)
		}

		return nil
	}

	// Deliver mail to local user inbox
	h, err := hash.NewFromHash(addrInfo.Hash)
	if err != nil {
		return err
	}

	// Add message
	ar := container.Instance.GetAccountRepo()
	err = ar.CreateMessage(*h, msgID)
	if err != nil {
		return err
	}

	// Move to inbox
	err = ar.AddToBox(*h, account.BoxInbox, msgID)
	if err != nil {
		return err
	}

	_ = dispatcher.DispatchLocalDelivery(*h, header, msgID)

	return nil
}

// deliverRemote uploads a message to a remote mail server. For this to work it first needs to fetch a
// ticket from that server. Either that ticket is supplied, or we need to do proof-of-work first before
// we get the ticket. Once we have the ticket, we can upload the message to the server in the same way
// we upload a message from a client to a server.
func deliverRemote(addrInfo *resolver.AddressInfo, msgID string, header *message.Header) error {
	rs := container.Instance.GetResolveService()

	_, err := message.GetRetryInfo(message.SectionProcessing, msgID)
	if err == nil {
		rs.ClearRoutingCacheEntry(addrInfo.RoutingID)
	}

	routingInfo, err := rs.ResolveRouting(addrInfo.RoutingID)
	if err != nil {
		logrus.Warnf("cannot find routing ID %s for %s. Retrying.", addrInfo.RoutingID, header.To.Addr)
		MoveToRetryQueue(msgID)
		return err
	}

	logrus.Debugf("Message %s is remote, transferring to %s", msgID, routingInfo.Routing)

	t, err := processTicket(*routingInfo, *addrInfo, header, msgID)
	if err != nil {
		return err
	}

	c, err := getClient(*routingInfo)
	if err != nil {
		logrus.Warning("cannot create API: ", err)
		return err
	}

	// parallelize uploads
	g := new(errgroup.Group)

	g.Go(func() error {
		logrus.Tracef("uploading header for ticket %s", t.ID)
		return c.UploadHeader(*t, header)
	})

	g.Go(func() error {
		catalogPath, err := message.GetPath(message.SectionProcessing, msgID, "catalog")
		if err != nil {
			return err
		}

		catalogData, err := ioutil.ReadFile(catalogPath)
		if err != nil {
			return err
		}

		logrus.Tracef("uploading catalog for ticket %s", t.ID)
		return c.UploadCatalog(*t, catalogData)
	})

	messageFiles, err := message.GetFiles(message.SectionProcessing, msgID)
	if err != nil {
		_ = c.DeleteMessage(*t)
		return err
	}

	for _, messageFile := range messageFiles {
		g.Go(uploadBlockFromFile(c, t, messageFile))
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		logrus.Debugf("Error while uploading message %s: %s", msgID, err)
		_ = c.DeleteMessage(*t)
		return err
	}

	// All done, mark upload as completed
	logrus.Tracef("message completed for ticket %s", t.ID)
	err = c.CompleteUpload(*t)
	if err != nil {
		return err
	}

	h, err := hash.NewFromHash(addrInfo.Hash)
	if err == nil {
		_ = dispatcher.DispatchRemoteDelivery(*h, header, msgID)
	}

	// Remove local message from processing queue
	return message.RemoveMessage(message.SectionProcessing, msgID)
}

func uploadBlockFromFile(c *api.API, t *ticket.Ticket, mf message.FileType) func() error {
	return func() error {
		// Open reader
		f, err := os.Open(mf.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		logrus.Tracef("uploading block %s for ticket %s", mf.ID, t.ID)
		return c.UploadBlock(*t, mf.ID, f)
	}
}

// processTicket will fetch a ticket from the mail server and validate it through proof-of-work
func processTicket(routingInfo resolver.RoutingInfo, addrInfo resolver.AddressInfo, header *message.Header, msgID string) (*ticket.Ticket, error) {
	// Get upload ticket
	h, err := hash.NewFromHash(addrInfo.Hash)
	if err != nil {
		return nil, err
	}

	logrus.Tracef("getting ticket for %s:%s:%s", header.From.Addr, *h, "")

	c, err := getClient(routingInfo)
	if err != nil {
		return nil, err
	}

	t, err := c.GetTicket(header.From.Addr, *h, "")
	if err != nil {
		return nil, err
	}

	// If the ticket is valid, then we are done
	if t.Valid {
		return t, nil
	}

	logrus.Debugf("ticket %s not valid. Need to do proof of work", t.ID)

	// Do proof of work. We have to wait for it. This is ok as this is just a separate thread.
	t.Work.Data.Work()

	logrus.Debugf("work for %s is completed", t.ID)
	t, err = c.ValidateTicket(header.From.Addr, header.To.Addr, t.SubscriptionID, t)
	if err != nil || !t.Valid {
		logrus.Warnf("Ticket for message %s not valid after proof of work, moving to retry queue", msgID)
		MoveToRetryQueue(msgID)
		return nil, errValidating
	}

	// TIcket is ok after we done our proof-of-work
	return t, nil
}

// getClient will return an API client pointing to the actual mail server found in the routing info
func getClient(routingInfo resolver.RoutingInfo) (*api.API, error) {
	return api.NewAnonymous(routingInfo.Routing, nil)
}
