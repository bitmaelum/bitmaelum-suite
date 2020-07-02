package processor

var (
	// UploadChannel Message has been uploaded
	UploadChannel chan string
	// OutgoingChannel Message is queued for outgoing server
	OutgoingChannel chan string
	// IncomingChannel Message is incoming from other server
	IncomingChannel chan string
)