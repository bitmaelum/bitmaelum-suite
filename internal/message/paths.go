package message

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"path/filepath"
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
	// SectionProcessing processes a message
	SectionProcessing
	// SectionRetry messages that have to be retried at a later stadium
	SectionRetry
)

// GetPath will return the actual path based on the section, messageID and file inside the message
func GetPath(section Section, msgID, file string) (string, error) {
	switch section {
	case SectionIncoming:
		return filepath.Join(config.Server.Paths.Incoming, msgID, file), nil
	case SectionRetry:
		return filepath.Join(config.Server.Paths.Retry, msgID, file), nil
	case SectionProcessing:
		return filepath.Join(config.Server.Paths.Processing, msgID, file), nil
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

	_, err = fs.Stat(p)
	return err == nil
}

// ProcessQueuePathExists returns true when the processing path for the given message/file exists
func ProcessQueuePathExists(msgID, file string) bool {
	p, err := GetPath(SectionProcessing, msgID, file)
	if err != nil {
		return false
	}

	_, err = fs.Stat(p)
	return err == nil
}
