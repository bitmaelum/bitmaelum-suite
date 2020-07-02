package message

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
)

const (
	// UploadPath Path where message are uploaded from a client
	UploadPath          string = ".upload"

	// ProcessQueuePath Path where messages are moved while processing for outgoing
	ProcessQueuePath    string = ".processing"

	// IncomingPath Path where incoming messages are stored
	IncomingPath        string = ".incoming"
)

func getUploadPath(uuid, file string) (string, error) {
	return homedir.Expand(path.Join(UploadPath, uuid, file))
}

func getProcessQueuePath(uuid, file string) (string, error) {
	return homedir.Expand(path.Join(ProcessQueuePath, uuid, file))
}

func getIncomingPath(uuid, file string) (string, error) {
	return homedir.Expand(path.Join(IncomingPath, uuid, file))
}

// UploadPathExists returns true when the upload path for the given message/file exists
func UploadPathExists(uuid, file string) bool {
	p, err := getUploadPath(uuid, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}

// ProcessQueuePathExists returns true when the processing path for the given message/file exists
func ProcessQueuePathExists(uuid, file string) bool {
	p, err := getProcessQueuePath(uuid, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}

// IncomingPathExists returns true when the incoming path for the given message/file exists
func IncomingPathExists(uuid, file string) bool {
	p, err := getIncomingPath(uuid, file)
	if err != nil {
		return false
	}

	_, err = os.Stat(p)
	return err == nil
}

