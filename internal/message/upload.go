package message

// Functions for message that are uploaded from clients

import (
	"encoding/json"
	"github.com/spf13/afero"
	"io"
	"path/filepath"
	"regexp"
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
