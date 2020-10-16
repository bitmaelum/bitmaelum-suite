// Copyright (c) 2020 BitMaelum Authors
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
	"time"

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

			// Message has been retried over 10 times. It's not gonna happen.
			logrus.Errorf("Message %s stuck in retry queue for too long. Giving up.", info.MsgID)
			err := message.RemoveMessage(message.SectionProcessing, info.MsgID)
			if err != nil {
				logrus.Warnf("Cannot remove message %s from the process queue.", info.MsgID)
				continue
			}
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
		info.LastRetriedAt = time.Now()
		info.RetryAt = time.Now().Add(getNextRetryDuration(info.Retries))
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
	return info.RetryAt.Unix() < time.Now().Unix()
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
