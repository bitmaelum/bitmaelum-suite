// Copyright (c) 2021 BitMaelum Authors
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

package imapgw

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/backend/backendutil"
	"github.com/sirupsen/logrus"
)

// Delimiter for mailboxes
var Delimiter = "/"

// Mailbox to hold messages and other info
type Mailbox struct {
	Subscribed bool
	Messages   []*Message

	name string
	user *User
	id   int
}

// Name of the mailbox
func (mbox *Mailbox) Name() string {
	return mbox.name
}

// Info from the mailbox
func (mbox *Mailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Delimiter: Delimiter,
		Name:      mbox.name,
	}
	return info, nil
}

func (mbox *Mailbox) uidNext() uint32 {
	var uid uint32
	for _, msg := range mbox.Messages {
		if msg.UID > uid {
			uid = msg.UID
		}
	}
	uid++
	return uid
}

func (mbox *Mailbox) flags() []string {
	flagsMap := make(map[string]bool)
	for _, msg := range mbox.Messages {
		for _, f := range msg.Flags {
			if !flagsMap[f] {
				flagsMap[f] = true
			}
		}
	}

	var flags []string
	for f := range flagsMap {
		flags = append(flags, f)
	}
	return flags
}

func (mbox *Mailbox) unseenSeqNum() uint32 {
	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		seen := false
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				seen = true
				break
			}
		}

		if !seen {
			return seqNum
		}
	}
	return 0
}

// Status of the mailbox
func (mbox *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(mbox.name, items)
	status.Flags = mbox.flags()
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = mbox.unseenSeqNum()

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = uint32(len(mbox.Messages))
		case imap.StatusUidNext:
			status.UidNext = mbox.uidNext()
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = 0 // TODO
		case imap.StatusUnseen:
			status.Unseen = 0 // TODO
		}
	}

	return status, nil
}

// SetSubscribed will mark the mailbox as subscribed
func (mbox *Mailbox) SetSubscribed(subscribed bool) error {
	mbox.Subscribed = subscribed
	return nil
}

// Check the mailbox
func (mbox *Mailbox) Check() error {
	return nil
}

// ListMessages of the mailbox
func (mbox *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	var err error
	mbox.Messages, err = refreshMailbox(mbox.user, mbox.id, mbox.Messages)
	if err != nil {
		return err
	}

	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		var id uint32
		if uid {
			id = msg.UID
		} else {
			id = seqNum
		}
		if !seqSet.Contains(id) {
			continue
		}

		m, err := msg.Fetch(seqNum, items, mbox.user)
		if err != nil {
			continue
		}

		ch <- m
	}

	return nil
}

// SearchMessages in the mailbox
func (mbox *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	var ids []uint32
	for i, msg := range mbox.Messages {
		seqNum := uint32(i + 1)

		ok, err := msg.Match(seqNum, criteria)
		if err != nil || !ok {
			continue
		}

		var id uint32
		if uid {
			id = msg.UID
		} else {
			id = seqNum
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateMessage will create a new local message
func (mbox *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return errUnimplemented
}

// UpdateMessagesFlags will update flags
func (mbox *Mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	for i, msg := range mbox.Messages {
		var id uint32
		if uid {
			id = msg.UID
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		msg.Flags = backendutil.UpdateFlags(msg.Flags, op, flags)

		logrus.Info(filepath.Join(os.TempDir(), "bm-bridge"))

		if _, err := os.Stat(filepath.Join(os.TempDir(), "bm-bridge")); os.IsNotExist(err) {
			logrus.Info("create tmpdir")
			os.Mkdir(filepath.Join(os.TempDir(), "bm-bridge"), os.ModeDir)
		}

		tmpfn := filepath.Join(os.TempDir(), "bm-bridge", msg.ID)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(msg.Flags)
		err := ioutil.WriteFile(tmpfn, buf.Bytes(), 0666)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

func getIDFromName(boxName string) int {
	switch boxName {
	case folderInbox:
		return 1
	case folderSent:
		return 2
	case folderTrash:
		return 3
	default:
		boxID, _ := strconv.Atoi(strings.TrimPrefix(boxName, "BOX_"))
		return boxID
	}
}

// CopyMessages will copy a message
func (mbox *Mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, destName string) error {
	boxID := getIDFromName(destName)
	if boxID == 0 {
		return backend.ErrNoSuchMailbox
	}

	for i, msg := range mbox.Messages {
		var id uint32
		if uid {
			id = msg.UID
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		mbox.user.Client.MoveMessage(mbox.user.Info.Address.Hash(), msg.ID, mbox.id, boxID)
	}

	return nil
}

// Expunge will delete a message
func (mbox *Mailbox) Expunge() error {
	for i := len(mbox.Messages) - 1; i >= 0; i-- {
		msg := mbox.Messages[i]

		deleted := false
		for _, flag := range msg.Flags {
			if flag == imap.DeletedFlag {
				err := mbox.user.Client.RemoveMessageFromBox(mbox.user.Info.Address.Hash(), msg.ID, mbox.id)
				if err != nil {
					return err
				}
				deleted = true
				break
			}
		}

		if deleted {
			mbox.Messages = append(mbox.Messages[:i], mbox.Messages[i+1:]...)
		}
	}

	return nil
}
