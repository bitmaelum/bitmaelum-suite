package processor

import "github.com/sirupsen/logrus"

var (
	// UploadChannel Message has been uploaded
	UploadChannel chan string
	// OutgoingChannel Message is queued for outgoing server
	OutgoingChannel chan string
	// IncomingChannel Message is incoming from other server
	IncomingChannel chan string
)

// QueueIncomingMessage queues message on the incoming channel
func QueueIncomingMessage(msgID string) {
	logrus.Tracef("queueing msgID message from incoming: %s", msgID)
	IncomingChannel <- msgID
}

// QueueOutgoingMessage queues message on the outgoing channel
func QueueOutgoingMessage(msgID string) {
	logrus.Tracef("queueing msgID message from upload: %s", msgID)
	OutgoingChannel <- msgID
}

// QueueUploadMessage queues message on the uploaded channel
func QueueUploadMessage(msgID string) {
	logrus.Tracef("queueing msgID message from upload: %s", msgID)
	UploadChannel <- msgID
}
