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
func ProcessRetryQueue() error {
	logrus.Debugf("scanning retry queue for action")

	retryQueue, err := message.GetRetryInfoFromQueue()
	if err != nil {
		return err
	}

	for _, info := range retryQueue {
		if info.Retries > MaxRetries {
			// @TODO: We should send a message back to the user?

			// Message has been retried over 10 times. It's not gonna happen.
			logrus.Errorf("Message %s stuck in retry queue for too long. Giving up.", info.MsgID)
			err := message.RemoveMessage(message.SectionProcessQueue, info.MsgID)
			if err != nil {
				logrus.Warnf("Cannot remove message %s from the process queue.", info.MsgID)
				continue
			}
		}

		if canRetryNow(info) {
			err := message.MoveMessage(message.SectionRetry, message.SectionProcessQueue, info.MsgID)
			if err != nil {
				continue
			}

			go ProcessMessage(info.MsgID)
		}
	}

	return nil
}

// MoveToRetryQueue moves a message (back) to retry queue
func MoveToRetryQueue(msgID string) {
	info, err := message.GetRetryInfo(message.SectionProcessQueue, msgID)
	if err == nil {
		info.Retries++
		info.LastRetriedAt = time.Now()
		info.RetryAt = time.Now().Add(getNextRetryDuration(info.Retries))
	} else {
		info = message.NewRetryInfo()
	}

	err = message.StoreRetryInfo(message.SectionProcessQueue, msgID, *info)
	if err != nil {
		logrus.Warnf("Cannot store retry information for message %s.", msgID)
	}

	err = message.MoveMessage(message.SectionProcessQueue, message.SectionRetry, info.MsgID)
	if err != nil {
		// can't move the message?
		logrus.Warnf("Cannot move message %s from processing to retry queue.", msgID)
	}
}

// canRetryNow returns true if we can retry the message right now
func canRetryNow(info message.RetryInfo) bool {
	if info.RetryAt.Unix() < time.Now().Unix() {
		return true
	}

	return false
}

// calculateNextRetryTime will return the next time a message can be retried again
func getNextRetryDuration(retries int) (d time.Duration) {
	// @TODO: These duration should be configurable. Something like this in config: "5:1 17:5 25:30 30:60"

	/*
	 * config:
	 *   retries: [
	 *     { count:  5, hold:  1 },
	 *     { count: 17, hold:  5 },
	 *     { count: 25, hold: 30 },
	 *     { count: 30, hold: 60 },
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
