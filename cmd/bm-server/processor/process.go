package processor

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/server"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"
)

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

	rs := container.GetResolveService()
	res, err := rs.ResolveAddress(header.To.Addr)
	if err != nil {
		logrus.Trace(err)
		logrus.Warnf("cannot resolve address %s for message %s. Retrying.", header.To.Addr, msgID)
		MoveToRetryQueue(msgID)
		return
	}

	// Local addresses don't need to be send. They are treated locally
	ar := container.GetAccountRepo()
	if ar.Exists(header.To.Addr) {
		// probably move the message to the incoming queue
		// Do stuff locally
		logrus.Debugf("Message %s can be transferred locally to %s", msgID, res.Hash)

		// Check the serverSignature
		if !server.VerifyHeader(*header) {
			logrus.Errorf("message %s destined for %s has failed the server signature check. Seems that this message did not originate from the original mail server. Removing the message.", msgID, header.To.Addr)

			err := message.RemoveMessage(message.SectionProcessing, msgID)
			if err != nil {
				logrus.Warnf("Cannot remove message %s from the process queue.", msgID)
			}
		}

		err := deliverLocal(res, msgID)
		if err != nil {
			logrus.Warnf("cannot deliver message %s locally to %s. Retrying.", msgID, header.To.Addr)
			MoveToRetryQueue(msgID)
		}
		return
	}

	routingRes, err := rs.ResolveRouting(res.RoutingID)
	if err != nil {
		logrus.Warnf("cannot find routing ID %s for %s. Retrying.", res.RoutingID, header.To.Addr)
		MoveToRetryQueue(msgID)
		return
	}

	// Otherwise, send to outgoing server
	logrus.Debugf("Message %s is remote, transferring to %s", msgID, routingRes.Routing)
	err = deliverRemote(header, res, routingRes, msgID)
	if err != nil {
		logrus.Warnf("cannot deliver message %s remotely to %s. Retrying.", msgID, header.To.Addr)
		MoveToRetryQueue(msgID)
	}
}

// deliverLocal moves a message to a local mailbox. This is an easy process as it only needs to move
// the message to another directory.
func deliverLocal(info *resolver.AddressInfo, msgID string) error {
	// Deliver mail to local user's inbox
	ar := container.GetAccountRepo()
	err := ar.SendToBox(address.HashAddress(info.Hash), account.BoxInbox, msgID)
	if err != nil {
		// Something went wrong.. let's try and move the message back to the retry queue
		logrus.Warnf("cannot deliver %s locally. Moving to retry queue", msgID)
		MoveToRetryQueue(msgID)
	}

	return nil
}

// deliverRemote uploads a message to a remote mail server. For this to work it first needs to fetch a
// ticket from that server. Either that ticket is supplied, or we need to do proof-of-work first before
// we get the ticket. Once we have the ticket, we can upload the message to the server in the same way
// we upload a message from a client to a server.
func deliverRemote(header *message.Header, info *resolver.AddressInfo, routingInfo *resolver.RoutingInfo, msgID string) error {
	client, err := api.NewAnonymous(api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Server.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
	if err != nil {
		return err
	}

	// Get upload ticket
	logrus.Tracef("getting ticket for %s:%s:%s", header.From.Addr, address.HashAddress(info.Hash), "")
	t, err := client.GetAnonymousTicket(header.From.Addr, address.HashAddress(info.Hash), "")
	if err != nil {
		return err
	}
	if !t.Valid {
		logrus.Debugf("ticket %s not valid. Need to do proof of work", t.ID)
		// Do proof of work. We have to wait for it. This is ok as this is just a separate thread.
		t.Proof.WorkMulticore()

		logrus.Debugf("work for %s is completed", t.ID)
		t, err = client.GetAnonymousTicketByProof(header.From.Addr, header.To.Addr, t.SubscriptionID, t.ID, t.Proof.Proof)
		if err != nil || !t.Valid {
			logrus.Warnf("Ticket for message %s not valid after proof of work, moving to retry queue", msgID)
			MoveToRetryQueue(msgID)
			return err
		}
	}

	// parallelize uploads
	g := new(errgroup.Group)
	g.Go(func() error {
		logrus.Tracef("uploading header for ticket %s", t.ID)
		return client.UploadHeader(*t, header)
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
		return client.UploadCatalog(*t, catalogData)
	})

	messageFiles, err := message.GetFiles(message.SectionProcessing, msgID)
	if err != nil {
		_ = client.DeleteMessage(*t)
		return err
	}

	for _, messageFile := range messageFiles {
		// Store locally, otherwise the anonymous go function doesn't know which "block"
		mf := messageFile

		g.Go(func() error {
			// Open reader
			f, err := os.Open(mf.Path)
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()

			logrus.Tracef("uploading block %s for ticket %s", mf.ID, t.ID)
			return client.UploadBlock(*t, mf.ID, f)
		})
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		logrus.Debugf("Error while uploading message %s: %s", msgID, err)
		_ = client.DeleteMessage(*t)
		return err
	}

	// All done, mark upload as completed
	logrus.Tracef("message completed for ticket %s", t.ID)
	err = client.CompleteUpload(*t)
	if err != nil {
		return err
	}

	// Remove local message from processing queue
	return message.RemoveMessage(message.SectionProcessing, msgID)
}
