package processor

import (
	"github.com/bitmaelum/bitmaelum-server/internal/message"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// ProcessMessage will process a message found in the processing queue.
//
//   * If it's a local address, it will be moved to the local mailbox
//   * If it's a remote address, it will be send to the remote mailserver
//   * If things fail, it will be moved to the retry queue, where it will be moved to processed queue later
func ProcessRetryQueue() error {
	logrus.Debugf("scanning retry queue for action")

	files, err := message.GetMessagesFromRetryQueue()
	if err != nil {
		return err
	}

	for _, f := range files {
		msgID := f.Name()

		needsRetry := checkForRetry(msgID)
		if needsRetry {
			message.MoveToProcessing(message.SectionRetry, msgID)
			go ProcessMessage(msgID)
		}
	}
}

