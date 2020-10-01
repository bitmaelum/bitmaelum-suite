package processor

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/sirupsen/logrus"
	"time"
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
	 *     { count:  5, hold:  1 },
	 *     { count: 17, hold:  5 },
	 *     { count: 25, hold: 30 },
	 *     { count: 30, hold: 60 }
	 *   ]
	 */

	d = 0

	switch {
	case retries < 5:
		d = 1 * time.Minute
	case retries < 17:
		d = 5 * time.Minute
	case retries < 25:
		d = 30 * time.Minute
	case retries < 30:
		d = 60 * time.Minute
	}

	return
}
