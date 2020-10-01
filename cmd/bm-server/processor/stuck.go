package processor

import (
	"io/ioutil"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/sirupsen/logrus"
)

// ProcessStuckIncomingMessages will process stuck message found in the incoming queue.
func ProcessStuckIncomingMessages() {
	p, err := message.GetPath(message.SectionIncoming, "", "")
	if err != nil {
		return
	}

	// Check all files in the directory
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}

	ticketRepo := container.GetTicketRepo()
	for _, fileInfo := range files {
		// Not a dir, so not a message
		if !fileInfo.IsDir() {
			continue
		}

		// Find corresponding ticket for this incoming message
		_, err := ticketRepo.Fetch(fileInfo.Name())
		if err != nil {
			logrus.Errorf("found message %s in incoming queue without ticket", fileInfo.Name())

			// @TODO Ticket / message is not found. What to do with the message?
		}

	}
}

// ProcessStuckProcessingMessages will process stuck message found in the processing queue.
func ProcessStuckProcessingMessages() {
	p, err := message.GetPath(message.SectionProcessing, "", "")
	if err != nil {
		return
	}

	// Check all files in the directory
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}

	for _, fileInfo := range files {
		// Not a dir, so not a message
		if !fileInfo.IsDir() {
			continue
		}

		// If the message is not in the scoreboard, we are not processing it at the moment, so we can move
		// it to the retry queue
		if !IsInScoreboard(message.SectionProcessing, fileInfo.Name()) {
			MoveToRetryQueue(fileInfo.Name())
		}
	}
}
