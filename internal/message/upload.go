package message

// Functions for message that are uploaded from clients

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// GetMessageHeader Returns a marshalled message header
func GetMessageHeader(section Section, uuid string) (*Header, error) {
	p, err := getPath(section, uuid, "header.json")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	header := &Header{}
	err = json.Unmarshal(data, &header)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// RemoveMessage removes a complete message (header, catalog, blocks etc)
func RemoveMessage(uuid string) error {
	p, err := getPath(SectionUpload, uuid, "")
	if err != nil {
		return err
	}

	return os.RemoveAll(p)
}

// StoreBlock stores a message block to disk
func StoreBlock(uuid, blockID string, r io.Reader) error {
	p, err := getPath(SectionUpload, uuid, blockID)
	if err != nil {
		return err
	}

	// Create path if needed
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to block file
	blockFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	_, err = io.Copy(blockFile, r)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		return err
	}

	return nil
}

// StoreCatalog stores a catalog to disk
func StoreCatalog(uuid string, r io.Reader) error {
	p, err := getPath(SectionUpload, uuid, "catalog")
	if err != nil {
		return err
	}

	// Create path if needed
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to catalog file
	catFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	_, err = io.Copy(catFile, r)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		return err
	}

	return nil
}

// StoreMessageHeader stores a message header to disk
func StoreMessageHeader(uuid string, header *Header) error {
	p, err := getPath(SectionUpload, uuid, "header.json")
	if err != nil {
		return err
	}

	// Create path if needed
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to catalog file
	headerFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}

	// Marshal data and save
	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	_, err = headerFile.Write(data)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		return err
	}

	return nil
}

// MoveToProcessing moves a message from incoming to processing
func MoveToProcessing(section Section, uuid string) error {
	oldPath, err := getPath(section, uuid, "")
	if err != nil {
		return err
	}

	newPath, err := getPath(SectionProcessQueue, uuid, "")
	if err != nil {
		return err
	}

	return os.Rename(oldPath, newPath)
}
