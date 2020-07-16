package message

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// RetryInfo is a structure that holds information about when a message has been retried or when it needs to be retried
type RetryInfo struct {
	RetryAt       time.Time `json:"retry_at"`        // Retry processing again on or after this time
	LastRetriedAt time.Time `json:"last_retried_at"` // Last time the message was processed/retried
	Retries       int       `json:"retries"`         // Number of retries already done
	MsgID         string    `json:"message_id"`      // Actual message ID (redundant since it's always inside the message directory)
}

// NewRetryInfo returns a new retry info structure
func NewRetryInfo(msgID string) *RetryInfo {
	return &RetryInfo{
		RetryAt:       time.Now().Add(60 * time.Second),
		LastRetriedAt: time.Now(),
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
	files, err := ioutil.ReadDir(p)
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

	data, err := ioutil.ReadFile(p)
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

	return ioutil.WriteFile(p, data, 0600)
}
