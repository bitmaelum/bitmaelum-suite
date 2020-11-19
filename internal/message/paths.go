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

package message

import (
	"errors"
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

/*
 * Finding paths on the mail server is a bit difficult. A message can be in different stages:
 *
 *   - A message can be uploaded by a client and be unprocessed yet.
 *   - A message can be inside the processing queue
 *   - A message can be inside the retry queue
 *   - A message can be just uploaded by another server (or locally) and waiting inside the incoming queue
 */

var errUnknownSection = errors.New("unknown section")

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
		return "", errUnknownSection
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
