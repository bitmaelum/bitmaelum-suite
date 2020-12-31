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

package account

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const messageDir string = "messages"

func (r *fileRepo) FetchMessageHeader(addr hash.Hash, messageID string) (*message.Header, error) {
	header := &message.Header{}
	err := r.fetchJSON(addr, filepath.Join(messageDir, messageID, "header.json"), header)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (r *fileRepo) FetchMessageCatalog(addr hash.Hash, messageID string) ([]byte, error) {
	catalog, err := r.fetch(addr, filepath.Join(messageDir, messageID, "catalog"))
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (r *fileRepo) FetchMessageBlock(addr hash.Hash, messageID, blockID string) ([]byte, error) {
	block, err := r.fetch(addr, filepath.Join(messageDir, messageID, blockID))
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (r *fileRepo) FetchMessageAttachment(addr hash.Hash, messageID, attachmentID string) (rdr io.ReadCloser, size int64, err error) {
	return r.fetchReader(addr, filepath.Join(messageDir, messageID, attachmentID))
}

// Query messages inside mailbox
func (r *fileRepo) FetchListFromBox(addr hash.Hash, box int, since time.Time, offset, limit int) (*MessageList, error) {
	var list = &MessageList{
		Meta: MetaType{
			Total:    0,
			Returned: 0,
			Limit:    limit,
			Offset:   offset,
		},
		Messages: []MessageType{},
	}

	// @TODO: We don't do anything with offset yet.

	files, err := afero.ReadDir(r.fs, r.getPath(addr, getBoxAsString(box)))
	if err != nil {
		return nil, err
	}

	logrus.Trace("Fetching dir: ")
	for _, f := range files {
		if f.Mode()&os.ModeSymlink == 0 && !f.IsDir() {
			logrus.Trace("not a dir: ", f.Name())
			continue
		}

		list.Meta.Total++
		if list.Meta.Returned >= limit {
			logrus.Trace("limit reached")
			break
		}

		header := &message.Header{}
		err := r.fetchJSON(addr, filepath.Join(getBoxAsString(box), f.Name(), "header.json"), header)
		if err != nil {
			logrus.Trace("cannot find header.json")
			continue
		}
		catalog, err := r.fetch(addr, filepath.Join(getBoxAsString(box), f.Name(), "catalog"))
		if err != nil {
			logrus.Trace("cannot find catalog")
			continue
		}

		// Skip files if we have an offset
		if offset > 0 {
			logrus.Trace("offset not yet reached")
			offset--
			continue
		}

		if !since.IsZero() && f.ModTime().Before(since) {
			logrus.Trace("before since")
			// Skip, because it's before our "since" query
			continue
		}

		logrus.Trace("adding to list: ", f.Name())
		list.Meta.Returned++
		list.Messages = append(list.Messages, MessageType{
			ID:      f.Name(),
			Header:  *header,
			Catalog: catalog,
		})
	}

	return list, nil
}

func (r *fileRepo) MoveToBox(addr hash.Hash, srcBox, dstBox int, msgID string) error {
	srcPath := r.getPath(addr, filepath.Join(getBoxAsString(srcBox), msgID))
	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(dstBox), msgID))

	return r.fs.Rename(srcPath, dstPath)
}

// Move message into the account system. THis is basically a bridge between the processing section, and the accounts.
// I'm not sure if this needs to be here
func (r *fileRepo) CreateMessage(addr hash.Hash, msgID string) error {
	srcPath, err := message.GetPath(message.SectionProcessing, msgID, "")
	if err != nil {
		return err
	}

	// Finally, move the message to our message directory
	dstPath := r.getPath(addr, filepath.Join(messageDir, msgID))
	return r.fs.Rename(srcPath, dstPath)
}

// RemoveMessage Removes a message complete from the account
func (r *fileRepo) RemoveMessage(addr hash.Hash, msgID string) error {
	p := r.getPath(addr, filepath.Join(messageDir, msgID))
	err := r.fs.RemoveAll(p)
	if err != nil {
		return err
	}

	// Remove any message references from boxes
	boxes, err := r.GetAllBoxes(addr)
	if err != nil {
		return err
	}
	for _, box := range boxes {
		_ = r.RemoveFromBox(addr, box.ID, msgID)
	}

	return nil
}

// AddToBox Symlinks the message to the box
func (r *fileRepo) AddToBox(addr hash.Hash, boxID int, msgID string) error {
	// Check if we can symlink
	symlink, ok := r.fs.(afero.Symlinker)
	if !ok {
		return errors.New("symlinking is not implemented")
	}

	srcPath := r.getPath(addr, filepath.Join(messageDir, msgID))
	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(boxID), msgID))

	return symlink.SymlinkIfPossible(srcPath, dstPath)
}

// RemoveFromBox Unlink/remove the message from the box
func (r *fileRepo) RemoveFromBox(addr hash.Hash, boxID int, msgID string) error {
	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(boxID), msgID))

	return r.fs.Remove(dstPath)
}
