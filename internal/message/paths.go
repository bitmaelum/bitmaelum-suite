package message

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
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

// Section of the path we want to
type Section int

const (
	// SectionIncoming uploading message
	SectionIncoming = iota
	// SectionProcessQueue processes a message
	SectionProcessQueue
	// SectionRetry messages that have to be retried at a later stadium
	SectionRetry
)

// GetPath will return the actual path based on the section, messageID and file inside the message
func GetPath(section Section, msgID, file string) (string, error) {
	switch section {
	case SectionIncoming:
		return homedir.Expand(path.Join(config.Server.Paths.Incoming, msgID, file))
	case SectionRetry:
		return homedir.Expand(path.Join(config.Server.Paths.Retry, msgID, file))
	case SectionProcessQueue:
		return homedir.Expand(path.Join(config.Server.Paths.Processing, msgID, file))
	default:
		return "", errors.New("unknown section")
	}
}

// IncomingPathExists returns true when the upload path for the given message/file exists
func IncomingPathExists(msgID, file string) bool {
	p, err := GetPath(SectionIncoming, msgID, file)
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
