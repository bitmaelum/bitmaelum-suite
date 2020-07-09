package processor

import (
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/core/resolve"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

// ProcessMessage will process a message found in the processing queue.
//
//   * If it's a local address, it will be moved to the local mailbox
//   * If it's a remote address, it will be send to the remote mail server
//   * If things fail, it will be moved to the retry queue, where it will be moved to processed queue later
//
func ProcessMessage(msgID string) {
	logrus.Debugf("processing message %s", msgID)

	// Check header and get recipient
	header, err := message.GetMessageHeader(message.SectionProcessing, msgID)
	if err != nil {
		// cannot read header.. Let's move to retry queue
		MoveToRetryQueue(msgID)
		return
	}

	rs := container.GetResolveService()
	res, err := rs.Resolve(header.To.Addr)
	if err != nil {
		logrus.Errorf("cannot resolve address %s for message %s. Retrying.", header.To.Addr, msgID)
		MoveToRetryQueue(msgID)
		return
	}

	// Local addresses don't need to be send. They are treated locally
	if res.IsLocal() {
		// probably move the message to the incoming queue
		// Do stuff locally
		logrus.Debugf("Message %s can be transferred locally to %s", msgID, res.Hash)

		err := deliverLocal(res, msgID)
		if err != nil {
			logrus.Errorf("cannot deliver message %s locally to %s. Retrying.", msgID, header.To.Addr)
			MoveToRetryQueue(msgID)
		}
		return
	}

	// Otherwise, send to outgoing server
	logrus.Debugf("Message %s is remote, transferring to %s", msgID, res.Server)
	err = deliverRemote(res, msgID)
	if err != nil {
		logrus.Errorf("cannot deliver message %s remotely to %s. Retrying.", msgID, header.To.Addr)
		MoveToRetryQueue(msgID)
	}
}

// deliverLocal moves a message to a local mailbox
func deliverLocal(info *resolve.Info, msgID string) error {
	// Deliver mail to local user's inbox
	as := container.GetAccountService()
	err := as.Deliver(msgID, address.HashAddress(info.Hash))
	if err != nil {
		// Something went wrong.. let's try and move the message back to the retry queue
		MoveToRetryQueue(msgID)
	}

	return nil
}

// deliverRemote uploads a message to a remote mail server
func deliverRemote(info *resolve.Info, msgID string) error {
	return nil
}
