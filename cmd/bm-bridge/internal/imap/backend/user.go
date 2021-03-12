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
	"strconv"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/emersion/go-imap/backend"
)

// User holds the user info
type User struct {
	Account  string
	Vault    *vault.Vault
	Info     *vault.AccountInfo
	Client   *api.API
	Database *Storable
}

// Username from account
func (u *User) Username() string {
	return strings.Replace(u.Account, "!", "", -1)
}

func getNameFromID(boxID int) string {
	switch boxID {
	case 1:
		return folderInbox
	case 2:
		return folderSent
	case 3:
		return folderTrash
	default:
		return "BOX_" + strconv.Itoa(boxID)
	}
}

// ListMailboxes for the user
func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	mbl, err := u.Client.GetMailboxList(u.Info.Address.Hash())
	if err != nil {
		panic(err)
	}

	for _, box := range mbl.Boxes {
		mailbox, err := u.GetMailbox(getNameFromID(box.ID))
		if err != nil {
			return nil, err
		}
		mailboxes = append(mailboxes, mailbox)
	}

	return
}

// GetMailbox from server
func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
	mbl, err := u.Client.GetMailboxList(u.Info.Address.Hash())
	if err != nil {
		panic(err)
	}

	for _, box := range mbl.Boxes {
		if name == getNameFromID(box.ID) {
			messages, err := refreshMailbox(u, box.ID, nil)
			if err != nil {
				return nil, err
			}

			mailbox := &Mailbox{
				name:     getNameFromID(box.ID),
				id:       box.ID,
				user:     u,
				Messages: messages,
			}

			return mailbox, nil
		}
	}

	err = errNoSuchMailbox
	return nil, err
}

// CreateMailbox is not implemented yet
func (u *User) CreateMailbox(name string) error {
	return errUnimplemented
}

// DeleteMailbox is not implemented yet
func (u *User) DeleteMailbox(name string) error {
	return errUnimplemented
}

// RenameMailbox is not implemented yet
func (u *User) RenameMailbox(existingName, newName string) error {
	return errUnimplemented

}

// Logout from server
func (u *User) Logout() error {
	return nil
}
