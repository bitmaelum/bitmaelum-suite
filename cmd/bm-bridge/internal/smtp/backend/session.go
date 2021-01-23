package smtpgw

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

const (
	errInvalidDomain = "you can only send to " + common.DefaultDomain
	errInvalidFrom   = "invalid from address, account not found on vault"
)

// A Session is returned after successful login.
type Session struct {
	Account string
	Vault   *vault.Vault
	Info    *vault.AccountInfo
	Client  *api.API
	From    *address.Address
	To      *address.Address
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	from = common.EmailToAddr(from)

	if s.Account == "" {
		var err error
		s.Info, s.Client, err = common.GetClientAndInfo(s.Vault, from)
		if err != nil {
			return err
		}
	}

	if from != s.Info.Address.String() {
		return errors.New(errInvalidFrom)
	}

	// Check from address
	fromAddr, err := address.NewAddress(from)
	if err != nil {
		return err
	}

	s.From = fromAddr

	return nil
}

func (s *Session) Rcpt(to string) error {
	if !strings.Contains(to, "@") {
		if !strings.HasSuffix(to, "!") {
			return errors.New(errInvalidDomain)
		}
	} else {
		if !strings.Contains(to, common.DefaultDomain) {
			return errors.New(errInvalidDomain)
		}
	}

	to = common.EmailToAddr(to)

	// Check to address
	toAddr, err := address.NewAddress(to)
	if err != nil {
		return err
	}

	s.To = toAddr

	return nil
}

func (s *Session) Data(r io.Reader) error {

	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return err
	}

	// Extract from name from headers
	fromName := ""
	tmpFrom := strings.Split(env.GetHeader("From"), "<")
	if len(tmpFrom) > 0 {
		fromName = tmpFrom[0]
	}

	// Fetch both sender and recipient info
	svc := container.Instance.GetResolveService()
	senderInfo, err := svc.ResolveAddress(s.From.Hash())
	if err != nil {
		return err
	}
	recipientInfo, err := svc.ResolveAddress(s.To.Hash())
	if err != nil {
		return err
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(s.From, nil, fromName, s.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(s.To, nil, &recipientInfo.PublicKey)

	// Get blocks
	blocks := make([]string, 0)
	blocks = append(blocks, "default,"+env.Text)
	if env.HTML != "" {
		blocks = append(blocks, "html,"+env.HTML)
	}
	for _, part := range env.OtherParts {
		blocks = append(blocks, part.ContentType+","+string(part.Content))
	}

	// Get attachments
	attachments := make([]string, 0)
	for _, attachment := range env.Attachments {
		// We write the attachments temporary to disk so we can use it later on message.Compose,
		// however this needs to be improved so we don't need to write them to disk
		fName := filepath.Join(os.TempDir(), attachment.FileName)
		err = ioutil.WriteFile(fName, attachment.Content, 0644)
		if err != nil {
			return err
		}

		defer os.Remove(fName)

		attachments = append(attachments, fName)
	}

	for _, inline := range env.Inlines {
		// We write the inlines temporary to disk so we can use it later on message.Compose,
		// however this needs to be improved so we don't need to write them to disk
		fName := filepath.Join(os.TempDir(), inline.FileName)
		err = ioutil.WriteFile(fName, inline.Content, 0644)
		if err != nil {
			return err
		}

		defer os.Remove(fName)

		blocks = append(blocks, inline.ContentType+",file:"+fName)
	}

	// Compose mail
	envelope, err := message.Compose(addressing, env.GetHeader("Subject"), blocks, attachments)
	if err != nil {
		return err
	}

	// Send mail
	client, err := api.NewAuthenticated(*s.From, s.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return err
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}

	/*for _, attachment := range attachments {
		os.Remove(attachment)
	}*/

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
