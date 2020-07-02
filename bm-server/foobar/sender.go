package foobar

import (
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/sirupsen/logrus"
)

// ProcessIncomingClientMessage processes a new message send from a client via a go function
func ProcessIncomingClientMessage(uuid string) {
	logrus.Debugf("processing incoming message %s", uuid)
	go process(uuid)
}

func process(uuid string) {
	// // Check if we are already processing the uuid
	// if ! checkAndSetProcessingList(uuid) {
	// 	return
	// }
	// // Remove from processing list when completed
	// defer unsetProcessingList(uuid)

	// Fetch header
	header, err := message.GetMessageHeader(uuid)
	if err != nil {
		return
	}

	// 1. Move message to processing area
	logrus.Debugf("moving message %s to processing queue", uuid)
	err = message.MoveIncomingMessageToProcessingQueue(uuid)
	if err != nil {
		logrus.Errorf("cannot move message %s to processing queue", uuid)
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


	logrus.Debugf("Server to send message to is %s ", res.Address)
	// 3. Communicate with server and send message
	//    3.1 If not able, move to retry queue
	//    3.2. After X time, move message to sender that the mail could not be send
	// 4. Remove message from processing area
}
