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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

func (r *fileRepo) FetchMessageHeader(addr hash.Hash, box int, messageID string) (*message.Header, error) {
	header := &message.Header{}
	err := r.fetchJSON(addr, filepath.Join(getBoxAsString(box), messageID, "header.json"), header)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (r *fileRepo) FetchMessageCatalog(addr hash.Hash, box int, messageID string) ([]byte, error) {
	catalog, err := r.fetch(addr, filepath.Join(getBoxAsString(box), messageID, "catalog"))
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

func (r *fileRepo) FetchMessageBlock(addr hash.Hash, box int, messageID, blockID string) ([]byte, error) {
	block, err := r.fetch(addr, filepath.Join(getBoxAsString(box), messageID, blockID))
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (r *fileRepo) FetchMessageAttachment(addr hash.Hash, box int, messageID, attachmentID string) (rdr io.ReadCloser, size int64, err error) {
	return r.fetchReader(addr, filepath.Join(getBoxAsString(box), messageID, attachmentID))
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

	files, err := ioutil.ReadDir(r.getPath(addr, getBoxAsString(box)))
	if err != nil {
		return nil, err
	}

	logrus.Trace("Fethcing dir: ")
	for _, f := range files {
		if !f.IsDir() {
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

	return os.Rename(srcPath, dstPath)
}

// Send a message to specific box
// @TODO: This is a bit difficult: this is actually a bridge between the processing engine and the account storage
// it assumes that both are using files. We must thus find a way to transfer a message from the processing to account
// without assumptions. This probably means reading the message in-memory or something, and we don't like that either.
// So we have to come up with a better way....
func (r *fileRepo) SendToBox(addr hash.Hash, box int, msgID string) error {
	srcPath, err := message.GetPath(message.SectionProcessing, msgID, "")
	if err != nil {
		return err
	}

	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(box), msgID))
	// // If we have the inbox, the message is prefixed with the current timestamp (UTC). This allows us
	// // sort on time locally and we can just fetch from a specific time (ie: fetch all messages since 20 minutes ago)
	// if box == "inbox" {
	// 	dstPath = r.getPath(addr, filepath.Join(box, fmt.Sprintf("%d-%s", time.Now().Unix(), msgID)))
	// }
	return os.Rename(srcPath, dstPath)
}
