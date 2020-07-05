package processor

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
	IncomingChannel <- uuid
}

// QueueOutgoingMessage queues message on the outgoing channel
func QueueOutgoingMessage(uuid string) {
	OutgoingChannel <- uuid
}

// QueueUploadMessage queues message on the uploaded channel
func QueueUploadMessage(uuid string) {
	UploadChannel <- uuid
}
