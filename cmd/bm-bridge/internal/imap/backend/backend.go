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
	"errors"
	"sort"
	"strconv"
	"time"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/mailbox"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"github.com/sirupsen/logrus"
)

var (
	errUnimplemented          = errors.New("not implemented")
	errNoSuchMailbox          = errors.New("no such mailbox")
	errIncorrectAddressFormat = errors.New("incorrect address format specified")
	errAccountNotFound        = errors.New("account not found in vault")
)

const (
	folderInbox   = "INBOX"
	folderTrash   = "Deleted Messages"
	folderArchive = "Archive"
	folderNotes   = "Notes"
	folderDrafts  = "Drafts"
	folderJunk    = "Junk"
	folderSent    = "Sent Messages"
)

// Backend will hold the Vault to use
type Backend struct {
	Vault    *vault.Vault
	Database *Storable
}

// Login is called when logging in to IMAP. Any password is valid as long as the account is found in the vault
func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	account := username + "!"

	addr, err := address.NewAddress(account)
	if err != nil {
		return nil, errIncorrectAddressFormat
	}
	if !be.Vault.HasAccount(*addr) {
		return nil, errAccountNotFound
	}

	user := &User{Account: account, Vault: be.Vault, Database: be.Database}

	user.Info, user.Client, err = common.GetClientAndInfo(be.Vault, account)
	if err != nil {
		return nil, err
	}

	logrus.Infof("IMAP: user %s logged in", username)

	return user, nil
}

// New will create a new Backend
func New(v *vault.Vault, d Storable) *Backend {
	return &Backend{
		Vault:    v,
		Database: &d,
	}
}

func refreshMailbox(u *User, boxid int, currentMessages []*Message) ([]*Message, error) {
	msgList, err := u.Client.GetMailboxMessages(u.Info.Address.Hash(), strconv.Itoa(boxid), time.Time{})
	if err != nil {
		return nil, err
	}

	var finalMsgList []api.MailboxMessagesMessage

	messages := []*Message{}
	if currentMessages != nil {
		messages = currentMessages
		for _, msg := range msgList.Messages {
			var found = false
			for _, m := range currentMessages {
				if m.ID == msg.ID {
					found = true
					break
				}
			}
			if found {
				continue
			}

			finalMsgList = append(finalMsgList, msg)
		}

	} else {
		finalMsgList = msgList.Messages
	}

	// Sort messages first
	msort := mailbox.NewMessageSort(u.Info.GetActiveKey().PrivKey, finalMsgList, mailbox.SortDate, false)
	sort.Sort(&msort)

	for i, msg := range finalMsgList {

		if currentMessages != nil {
			i = i + len(currentMessages)
		}

		message := Message{
			UID:   uint32(i),
			ID:    msg.ID,
			User:  u,
			Flags: (*u.Database).Retrieve(msg.ID),
		}

		messages = append(messages, &message)
	}

	return messages, nil
}
