package processor

import "github.com/sirupsen/logrus"

var (
	// IncomingChannel Message with given msgID is incoming from a client
	IncomingChannel chan string
)

// QueueIncomingMessage queues message on the uploaded channel so it can be picked up for processing by the main loop
func QueueIncomingMessage(msgID string) {
	logrus.Tracef("queueing incoming message %s", msgID)
	IncomingChannel <- msgID
}
