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
	"fmt"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/sirupsen/logrus"
)

const (
	// MaxRetries defines how many retries we can do for sending a message
	MaxRetries int = 30
)

// ProcessRetryQueue will process all mails found in the retry queue or removes them when they are expired
func ProcessRetryQueue(forceRetry bool) {
	// Get retry info from all messages found in the retry queue
	retryQueue, err := message.GetRetryInfoFromQueue()
	if err != nil {
		return
	}

	for _, info := range retryQueue {
		if info.Retries > MaxRetries {
			// @TODO: We should send a message back to the user?
			logrus.Trace("Send a postmaster message")

			// get the message header since we need the from address
			msgHeader, _ := message.GetMessageHeader(message.SectionRetry, info.MsgID)

			logrus.Trace(msgHeader.From.Addr)
			_ = sendPostmasterMail(msgHeader.From.Addr, "Unable to deliver message", "The message with ID "+info.MsgID+" destinated to "+msgHeader.To.Addr.String()+" was unable to be delivered.")

			// Message has been retried over 10 times. It's not gonna happen.
			logrus.Errorf("Message %s stuck in retry queue for too long. Giving up.", info.MsgID)
			err = message.RemoveMessage(message.SectionRetry, info.MsgID)
			if err != nil {
				logrus.Warnf("Cannot remove message %s from the process queue.", info.MsgID)
			}
			continue
		}

		if forceRetry || canRetryNow(info) {
			err := message.MoveMessage(message.SectionRetry, message.SectionProcessing, info.MsgID)
			if err != nil {
				continue
			}

			go ProcessMessage(info.MsgID)
		}
	}
}

// MoveToRetryQueue moves a message (back) to retry queue and update retry info
func MoveToRetryQueue(msgID string) {
	// Create or update retry information for this message
	info, err := message.GetRetryInfo(message.SectionProcessing, msgID)
	if err == nil {
		info.Retries++
		info.LastRetriedAt = internal.TimeNow()
		info.RetryAt = internal.TimeNow().Add(getNextRetryDuration(info.Retries))
	} else {
		info = message.NewRetryInfo(msgID)
	}

	err = message.StoreRetryInfo(message.SectionProcessing, msgID, *info)
	if err != nil {
		logrus.Warnf("Cannot store retry information for message %s.", msgID)
	}

	// Move the message to the retry queue
	err = message.MoveMessage(message.SectionProcessing, message.SectionRetry, info.MsgID)
	if err != nil {
		// can't move the message?
		logrus.Warnf("Cannot move message %s from processing to retry queue.", msgID)
	}
}

// canRetryNow returns true if we can retry the message right now
func canRetryNow(info message.RetryInfo) bool {
	return info.RetryAt.Unix() <= internal.TimeNow().Unix()
}

// calculateNextRetryTime will return the next time a message can be retried again
func getNextRetryDuration(retries int) (d time.Duration) {
	/* @TODO: These duration should be configurable:
	 *
	 * config:
	 *   retries: [
	 *     { min:  0, max:  5, hold:  1 },
	 *     { min:  6, max: 17, hold:  5 },
	 *     { min: 18, max: 25, hold: 30 },
	 *     { min: 26, max: 0, hold: 60 }
	 *   ]
	 */

	switch {
	case retries <= 5:
		d = 1 * time.Minute
	case retries >= 6 && retries <= 17:
		d = 5 * time.Minute
	case retries >= 18 && retries <= 25:
		d = 30 * time.Minute
	case retries >= 26:
		d = 60 * time.Minute
	}

	return
}

func sendPostmasterMail(toHash hash.Hash, subject string, body string) error {
	// generate a hash from the server routingID to use it as from address
	fromHash, err := hash.NewFromHash(config.Routing.RoutingID)
	if err != nil {
		logrus.Trace("postmastermail: incorrect hash: ", err)
		return err
	}

	// Fetch public key from routing
	rs := container.Instance.GetResolveService()
	addr, err := rs.ResolveAddress(toHash)
	if err != nil {
		logrus.Trace("Unable to resolve address : ", err)
		return err
	}

	// Setup addressing
	senderName := fmt.Sprintf("postmaster at %s", config.Server.Server.Hostname)

	addressing := message.NewAddressing(message.SignedByTypeServer)
	addressing.AddSender(nil, fromHash, senderName, config.Routing.PrivateKey, "host")
	addressing.AddRecipient(nil, &toHash, &addr.PublicKey)

	// Create a single block with our body
	blocks := []string{"default," + body}

	logrus.Trace(addressing)

	envelope, err := message.Compose(addressing, subject, blocks, nil)
	if err != nil {
		logrus.Trace("Unable to create postmaster envelope : ", err)
		return err
	}

	// store message locally
	msgID, err := message.StoreLocalMessage(envelope)
	if err != nil {
		logrus.Trace("Unable to store local message : ", err)
		return err
	}

	// queue the message
	QueueIncomingMessage(msgID)

	return nil
}
