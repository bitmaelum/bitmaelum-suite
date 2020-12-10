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

// Functions for message that are uploaded from clients

import (
	"encoding/json"
	"io"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

var (
	uuidv4Regex = regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}")
)

// FileType is a simple message-id => path combination
type FileType struct {
	ID   string
	Path string
}

// GetMessageHeader Returns a marshalled message header
func GetMessageHeader(section Section, msgID string) (*Header, error) {
	p, err := GetPath(section, msgID, "header.json")
	if err != nil {
		return nil, err
	}

	data, err := afero.ReadFile(fs, p)
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

// GetFiles returns all blocks and attachments for the given message ID
func GetFiles(section Section, msgID string) ([]FileType, error) {
	p, err := GetPath(section, msgID, "")
	if err != nil {
		return nil, err
	}

	files, err := afero.ReadDir(fs, p)
	if err != nil {
		return nil, err
	}

	var ret []FileType

	for _, fi := range files {
		// skip dirs, "header.json" and "catalog"
		if fi.IsDir() || fi.Name() == "header.json" || fi.Name() == "catalog" {
			continue
		}

		// Only accept UUIDv4 filenames
		if !uuidv4Regex.MatchString(fi.Name()) {
			continue
		}

		ret = append(ret, FileType{
			ID:   fi.Name(),
			Path: filepath.Join(p, fi.Name()),
		})
	}

	return ret, nil
}

// RemoveMessage removes a complete message (header, catalog, blocks etc)
func RemoveMessage(section Section, msgID string) error {
	p, err := GetPath(section, msgID, "")
	if err != nil {
		return err
	}

	return fs.RemoveAll(p)
}

// StoreBlock stores a message block to disk
func StoreBlock(msgID, blockID string, r io.Reader) error {
	p, err := GetPath(SectionIncoming, msgID, blockID)
	if err != nil {
		return err
	}

	// Create path if needed
	err = fs.MkdirAll(filepath.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to block file
	blockFile, err := fs.Create(p)
	if err != nil {
		return err
	}
	defer func() {
		_ = blockFile.Close()
	}()

	_, err = io.Copy(blockFile, r)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = fs.Remove(p)
		return err
	}

	return nil
}

// StoreAttachment stores a message attachment to disk
func StoreAttachment(msgID, attachmentID string, r io.Reader) error {
	p, err := GetPath(SectionIncoming, msgID, attachmentID)
	if err != nil {
		return err
	}

	// Create path if needed
	err = fs.MkdirAll(filepath.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to attachment file
	attachmentFile, err := fs.Create(p)
	if err != nil {
		return err
	}
	defer func() {
		_ = attachmentFile.Close()
	}()

	_, err = io.Copy(attachmentFile, r)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = fs.Remove(p)
		return err
	}

	return nil
}

// StoreCatalog stores a catalog to disk
func StoreCatalog(msgID string, r io.Reader) error {
	p, err := GetPath(SectionIncoming, msgID, "catalog")
	if err != nil {
		return err
	}

	// Create path if needed
	err = fs.MkdirAll(filepath.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to catalog file
	catFile, err := fs.Create(p)
	if err != nil {
		return err
	}

	defer func() {
		_ = catFile.Close()
	}()

	_, err = io.Copy(catFile, r)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = fs.Remove(p)
		return err
	}

	return nil
}

// StoreHeader stores a message header to disk
func StoreHeader(msgID string, header *Header) error {
	p, err := GetPath(SectionIncoming, msgID, "header.json")
	if err != nil {
		return err
	}

	// Create path if needed
	err = fs.MkdirAll(filepath.Dir(p), 0777)
	if err != nil {
		return err
	}

	// Copy body straight to catalog file
	headerFile, err := fs.Create(p)
	if err != nil {
		return err
	}

	defer func() {
		_ = headerFile.Close()
	}()

	// Marshal data and save
	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	_, err = headerFile.Write(data)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = fs.Remove(p)
		return err
	}

	return nil
}

// MoveMessage moves a message from a section to another section. Highly unoptimized.
func MoveMessage(srcSection Section, targetSection Section, msgID string) error {
	p1, err := GetPath(srcSection, msgID, "")
	if err != nil {
		return err
	}

	// return if source path does not exist
	if _, err = fs.Stat(p1); err != nil {
		return err
	}

	// Create target path directories (if needed)
	p2, _ := GetPath(targetSection, msgID, "")
	err = fs.MkdirAll(filepath.Dir(p2), 0755)
	if err != nil {
		return err
	}

	return fs.Rename(p1, p2)
}

// StoreLocalMessage will store a message locally
func StoreLocalMessage(header *Header, catalog io.Reader, blocks map[string]*io.Reader, attachments map[string]*io.Reader) (string, error) {
	// create a uuid
	msgID, _ := uuid.NewRandom()

	// store the header
	err := StoreHeader(msgID.String(), header)
	if err != nil {
		return "", err
	}

	// store the catalog
	err = StoreCatalog(msgID.String(), catalog)
	if err != nil {
		return "", err
	}

	// store the blocks
	for id, r := range blocks {
		err = StoreBlock(msgID.String(), id, *r)
		if err != nil {
			return "", err
		}
	}

	// store the attachments
	for id, r := range attachments {
		err = StoreAttachment(msgID.String(), id, *r)
		if err != nil {
			return "", err
		}
	}

	// return the msgID
	return msgID.String(), nil
}
