package message

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

/*
 * Finding paths on the mail server is a bit difficult. A message can be in different stages:
 *
 *   - A message can be uploaded by a client and be unprocessed yet.
 *   - A message can be inside the processing queue
 *   - A message can be inside the retry queue
 *   - A message can be just uploaded by another server (or locally) and waiting inside the incoming queue
 */

const (
	// UploadPath Path where message are uploaded from a client
	UploadPath string = ".upload"

	// ProcessQueuePath Path where messages are moved while processing for outgoing
	ProcessQueuePath string = ".processing"

	// RetryPath Path where messages are stored before retrying again
	RetryPath string = ".retry"

	// IncomingPath Path where incoming messages are stored
	IncomingPath string = ".incoming"
)

// Section of the path we want to
type Section int

const (
	// SectionUpload uploading message
	SectionUpload = iota
	// SectionProcessQueue processes a message
	SectionProcessQueue
	// SectionRetry messages that have to be retried at a later stadium
	SectionRetry
	// SectionIncoming incoming message from other mailserver
	SectionIncoming
)

// GetPath will return the actual path based on the section, messageID and file inside the message
func GetPath(section Section, msgID, file string) (string, error) {
	switch section {
	case SectionUpload:
		return homedir.Expand(path.Join(UploadPath, msgID, file))
	case SectionRetry:
		return homedir.Expand(path.Join(RetryPath, msgID, file))
	case SectionProcessQueue:
		return homedir.Expand(path.Join(ProcessQueuePath, msgID, file))
	case SectionIncoming:
		return homedir.Expand(path.Join(IncomingPath, msgID, file))
	default:
		return "", errors.New("unknown section")
	}
}

// UploadPathExists returns true when the upload path for the given message/file exists
func UploadPathExists(msgID, file string) bool {
	p, err := GetPath(SectionUpload, msgID, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}

// ProcessQueuePathExists returns true when the processing path for the given message/file exists
func ProcessQueuePathExists(msgID, file string) bool {
	p, err := GetPath(SectionProcessQueue, msgID, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}

// IncomingPathExists returns true when the incoming path for the given message/file exists
func IncomingPathExists(msgID, file string) bool {
	p, err := GetPath(SectionIncoming, msgID, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}
