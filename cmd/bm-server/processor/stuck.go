// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

	ticketRepo := container.Instance.GetTicketRepo()
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
