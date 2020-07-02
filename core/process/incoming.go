package process

// Functions for message that are uploaded from clients

import (
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/sirupsen/logrus"
	"time"
)

// IncomingClientMessage processes a new message send from a client via a go function
func IncomingClientMessage(uuid string) {
	logrus.Debugf("processing incoming message %s", uuid)
	go processIncoming(uuid)
}

func processIncoming(uuid string) {
	// 1. Move message to processing area
	logrus.Debugf("moving message %s to processing queue", uuid)
	err = message.MoveIncomingMessageToProcessingQueue(uuid)
	if err != nil {
		logrus.Errorf("cannot move message %s to processing queue", uuid)
		return
	}

	// Fetch header
	header, err := message.GetMessageHeader(message.ProcessQueuePath, uuid)
	if err != nil {
		return
	}

	// 2. Check header for address and get server
	logrus.Debugf("resolving info for %s", header.To.Addr)

	rs := container.GetResolveService()
	res, err := rs.Resolve(header.To.Addr)
	if err != nil {
		logrus.Errorf("cannot resolve address %s for message %s", header.To.Addr, uuid)
		return
	}

	// Local addresses don't need to be send. They are treated locally
	if res.IsLocal() {
		// probably move the message to the incoming queue
		// Do stuff locally
		logrus.Debugf("Message %s can be transferred locally to %s", uuid, res.Hash)
	}

	// 3. Communicate with server and send message
	logrus.Debugf("Server to send message to is %s ", res.Address)
	err = ServerUpload(uuid, res.Address)
	if err != nil {
		// Schedule retry because we could not send the message
		ScheduleRetry(uuid)
	}

	// 4. Remove message from processing area
}
