package processor

import (
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/sirupsen/logrus"
)

// ProcessMessage will process a message found in the processing queue.
//
//   * If it's a local address, it will be moved to the local mailbox
//   * If it's a remote address, it will be send to the remote mail server
//   * If things fail, it will be moved to the retry queue, where it will be moved to processed queue later
func ProcessMessage(msgID string) {
	logrus.Debugf("processing message %s", msgID)

	// Check header
	header, err := message.GetMessageHeader(message.SectionProcessQueue, msgID)
	if err != nil {
		// cannot read header.. Let's move to retry queue
		MoveToRetryQueue(msgID)
		return
	}

	rs := container.GetResolveService()
	res, err := rs.Resolve(header.To.Addr)
	if err != nil {
		logrus.Errorf("cannot resolve address %s for message %s", header.To.Addr, msgID)
		return
	}

	// Local addresses don't need to be send. They are treated locally
	if res.IsLocal() {
		// probably move the message to the incoming queue
		// Do stuff locally
		logrus.Debugf("Message %s can be transferred locally to %s", msgID, res.Hash)

		return
	}

	// Otherwise, send to outgoing server
	logrus.Debugf("Message %s is remove, transferring to %s", msgID, res.Server)

}
