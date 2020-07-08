package processor

import "github.com/sirupsen/logrus"

var (
	// IncomingChannel Message is incoming from other server
	IncomingChannel chan string
)

// QueueIncomingMessage queues message on the uploaded channel
func QueueIncomingMessage(msgID string) {
	logrus.Tracef("queueing incoming message %s", msgID)
	IncomingChannel <- msgID
}
