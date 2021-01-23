package imapgw

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/emersion/go-imap/backend"
)

type User struct {
	Account string
	Vault   *vault.Vault
	Info    *vault.AccountInfo
	Client  *api.API
}

func (u *User) Username() string {
	return strings.Replace(u.Account, "!", "", -1)
}

func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	mbl, err := u.Client.GetMailboxList(u.Info.Address.Hash())
	if err != nil {
		panic(err)
	}

	for _, box := range mbl.Boxes {
		boxName := strconv.Itoa(box.ID)
		switch box.ID {
		case 1:
			boxName = "INBOX"
		case 2:
			boxName = "SENT"
		case 3:
			boxName = "TRASH"
		default:
			boxName = "BOX_" + boxName
		}
		mailbox, err := u.GetMailbox(boxName)
		if err != nil {
			return nil, err
		}
		mailboxes = append(mailboxes, mailbox)
	}

	return
}

func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
	mbl, err := u.Client.GetMailboxList(u.Info.Address.Hash())
	if err != nil {
		panic(err)
	}

	for _, box := range mbl.Boxes {
		boxName := strconv.Itoa(box.ID)
		switch box.ID {
		case 1:
			boxName = "INBOX"
		case 2:
			boxName = "SENT"
		case 3:
			boxName = "TRASH"
		default:
			boxName = "BOX_" + boxName
		}
		if name == boxName {
			messages, err := refreshMailbox(u, box.ID, nil)
			if err != nil {
				return nil, err
			}

			mailbox := &Mailbox{
				name:     boxName,
				id:       box.ID,
				user:     u,
				Messages: messages,
			}

			return mailbox, nil
		}
	}

	err = errors.New("No such mailbox")
	return nil, err
}

func (u *User) CreateMailbox(name string) error {
	return errors.New("Not implemented")
}

func (u *User) DeleteMailbox(name string) error {
	return errors.New("Not implemented")
}

func (u *User) RenameMailbox(existingName, newName string) error {
	return errors.New("Not implemented")

}

func (u *User) Logout() error {
	return nil
}
