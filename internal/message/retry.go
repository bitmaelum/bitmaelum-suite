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
	"encoding/json"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/spf13/afero"
)

// RetryInfo is a structure that holds information about when a message has been retried or when it needs to be retried
type RetryInfo struct {
	RetryAt       time.Time `json:"retry_at"`        // Retry processing again on or after this time
	LastRetriedAt time.Time `json:"last_retried_at"` // Last time the message was processed/retried
	Retries       int       `json:"retries"`         // Number of retries already done
	MsgID         string    `json:"message_id"`      // Actual message ID (redundant since it's always inside the message directory)
}

// override for testing purposes
var timeNow = internal.TimeNow

// NewRetryInfo returns a new retry info structure
func NewRetryInfo(msgID string) *RetryInfo {
	return &RetryInfo{
		RetryAt:       timeNow().Add(60 * time.Second),
		LastRetriedAt: timeNow(),
		Retries:       0,
		MsgID:         msgID,
	}
}

// GetRetryInfoFromQueue retrieves a list of retry infos as found in the retry queue
func GetRetryInfoFromQueue() ([]RetryInfo, error) {
	p, err := GetPath(SectionRetry, "", "")
	if err != nil {
		return []RetryInfo{}, err
	}

	// Check all files in the directory
	files, err := afero.ReadDir(fs, p)
	if err != nil {
		return []RetryInfo{}, err
	}

	var results []RetryInfo
	for _, fileInfo := range files {
		// Not a dir, so not a message
		if !fileInfo.IsDir() {
			continue
		}

		// Collect message retry info from message
		mri, err := GetRetryInfo(SectionRetry, fileInfo.Name())
		if err != nil {
			continue
		}

		results = append(results, *mri)
	}

	return results, nil
}

// GetRetryInfo will return information found in the message .retry.json file
func GetRetryInfo(section Section, msgID string) (*RetryInfo, error) {
	p, err := GetPath(section, msgID, ".retry.json")
	if err != nil {
		return nil, err
	}

	data, err := afero.ReadFile(fs, p)
	if err != nil {
		return nil, err
	}

	mri := &RetryInfo{}
	err = json.Unmarshal(data, &mri)
	if err != nil {
		return nil, err
	}

	return mri, err
}

// StoreRetryInfo saves the retry information back to disk
func StoreRetryInfo(section Section, msgID string, info RetryInfo) error {
	p, err := GetPath(section, msgID, ".retry.json")
	if err != nil {
		return err
	}

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, p, data, 0600)
}
