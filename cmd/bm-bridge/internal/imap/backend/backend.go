package imapgw

import (
	"errors"
	"sort"
	"strconv"
	"time"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/mailbox"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	Vault *vault.Vault
}

func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	account := username + "!"

	addr, err := address.NewAddress(account)
	if err != nil {
		return nil, errors.New("NO incorrect address format specified")
	}
	if !be.Vault.HasAccount(*addr) {
		return nil, errors.New("NO account not found in vault")
	}

	user := &User{Account: account, Vault: be.Vault}

	user.Info, user.Client, err = common.GetClientAndInfo(be.Vault, account)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func New(v *vault.Vault) *Backend {
	return &Backend{
		Vault: v,
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

		// Decrypt message if possible
		em := message.EncryptedMessage{
			ID:      msg.ID,
			Header:  &msg.Header,
			Catalog: msg.Catalog,

			GenerateBlockReader:      u.Client.GenerateAPIBlockReader(u.Info.Address.Hash()),
			GenerateAttachmentReader: u.Client.GenerateAPIAttachmentReader(u.Info.Address.Hash()),
		}

		_, err := em.Decrypt(u.Info.GetActiveKey().PrivKey)
		if err != nil {
			continue
		}

		message := Message{
			Uid:  uint32(i),
			ID:   msg.ID,
			User: u,
		}

		messages = append(messages, &message)
	}

	return messages, nil
}
