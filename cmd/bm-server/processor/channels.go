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
func QueueIncomingMessage(uuid string) {
	logrus.Tracef("queueing uuid message from incoming: %s", uuid)
	IncomingChannel <- uuid
}

// QueueOutgoingMessage queues message on the outgoing channel
func QueueOutgoingMessage(uuid string) {
	logrus.Tracef("queueing uuid message from upload: %s", uuid)
	OutgoingChannel <- uuid
}

// QueueUploadMessage queues message on the uploaded channel
func QueueUploadMessage(uuid string) {
	logrus.Tracef("queueing uuid message from upload: %s", uuid)
	UploadChannel <- uuid
}
