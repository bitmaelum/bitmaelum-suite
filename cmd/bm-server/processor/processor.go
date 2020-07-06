package processor

import "github.com/sirupsen/logrus"

// ProcessMessage will process a message found in the processing queue.
//
//   * If it's a local address, it will be moved to the local mailbox
//   * If it's a remote address, it will be send to the remote mailserver
//   * If things fail, it will be moved to the retry queue, where it will be moved to processed queue later
func ProcessMessage(msgID string) error {
	logrus.Debugf("processing message %s", msgID)

	// There is a message in the processing queue which we need to process.



}
