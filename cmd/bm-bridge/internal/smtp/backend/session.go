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

package smtpgw

import (
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-smtp"
	"github.com/mileusna/spf"
	"github.com/sirupsen/logrus"
)

var (
	errInvalidDomain  = errors.New("you can only send to @" + config.Bridge.Server.SMTP.Domain)
	errInvalidFrom    = errors.New("invalid from address, account not found on vault")
	errInvalidAddress = errors.New("invalid email address")
	errSPFnotPass     = errors.New("the SPF record not passed validation, are you a spammer? your days are over [https://github.com/bitmaelum/bitmaelum-suite/wiki/How-does-it-differ]")
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// A Session is returned after successful login.
type Session struct {
	Account    string
	Vault      *vault.Vault
	Info       *vault.AccountInfo
	Client     *api.API
	From       string
	IsGateway  bool
	To         string
	RemoteAddr net.Addr
}

// Mail is called when doing a "MAIL FROM:"
func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	// If running in gateway mode then it will check that the sender
	// is from outside bitmaelum network
	if s.IsGateway {
		if !isEmailValid(from) {
			logrus.Tracef("SMTP: discarding incoming mail from %s, invalid address", from)
			return errInvalidAddress
		}

		// Check that the mails comes from outside
		if isEmailFromBitMaelumNetwork(s.From) {
			logrus.Tracef("SMTP: discarding incoming mail from %s, the email is coming from inside", from)
			return nil
		}

		s.From = from
		return nil
	}

	from = common.EmailToAddr(from)

	if s.Account == "" {
		var err error
		s.Info, s.Client, err = common.GetClientAndInfo(s.Vault, from)
		if err != nil {
			return err
		}
	}

	if from != s.Info.Address.String() {
		return errInvalidFrom
	}

	s.From = from

	return nil
}

// Rcpt is called when doing a "RCPT TO:"
func (s *Session) Rcpt(to string) error {
	// If running in gateway mode then it will check that the receipient
	// belongs to bitmaelum network
	if s.IsGateway {
		// Check that the mails goes to @bitmaelum.network
		if strings.HasSuffix(to, "@"+config.Bridge.Server.SMTP.Domain) {
			s.To = common.EmailToAddr(to)
			return nil
		}

		logrus.Tracef("SMTP: discarding incoming mail to %s, address is not for ", config.Bridge.Server.SMTP.Domain)
		return errInvalidAddress
	}

	// Check if the recipient is in email format or bitmaelum format
	if !strings.Contains(to, "@") {
		if !strings.HasSuffix(to, "!") {
			return errInvalidDomain
		}
	} else {
		if !strings.HasSuffix(to, "@"+config.Bridge.Server.SMTP.Domain) {
			// If recipient is outside bitmaelum network then send the message to
			// the default gateway address
			if !isEmailValid(to) {
				return errInvalidAddress
			}

			if config.Bridge.Server.SMTP.GatewayAccount == "" {
				s.To = common.DefaultGatewayAccount
			} else {
				s.To = config.Bridge.Server.SMTP.GatewayAccount
			}

			return nil
		}
	}

	s.To = common.EmailToAddr(to)

	return nil
}

// Data is called when doing a "DATA"
func (s *Session) Data(r io.Reader) error {
	if s.IsGateway {
		// Since mail is from outside, at least check SPF record
		parts := strings.Split(s.From, "@")
		res := spf.CheckHost(net.ParseIP(strings.Split(s.RemoteAddr.String(), ":")[0]), parts[1], s.From, "")
		if res != spf.Pass {
			logrus.Tracef("SMTP: discarding incoming mail from %s, SPF did not pass", s.From)
			return errSPFnotPass
		}
	}

	logrus.Infof("SMTP: processing incoming email. From: %s, To: %s", s.From, s.To)

	return s.sendTo(r)
}

func isDomainHostedOnBitMaelumNetwork(domainName string) bool {
	// @TODO: We should probably check the DNS to check somehow if the domain is being actually hosted in the BitMaelum network
	if domainName == config.Bridge.Server.SMTP.Domain || domainName == common.DefaultDomain {
		return true
	}
	return false
}

func isEmailFromBitMaelumNetwork(emailAddress string) bool {
	atIndex := strings.LastIndex(emailAddress, "@")
	return isDomainHostedOnBitMaelumNetwork(emailAddress[atIndex+1:])
}

func isGatewayDestination(t string) bool {
	switch t {
	case config.Bridge.Server.SMTP.GatewayAccount:
		return true
	case common.DefaultGatewayAccount:
		return true
	}

	return false
}

// sendTo sends the mail to the bitmaelum network
func (s *Session) sendTo(r io.Reader) error {
	// Set up blocks & attachments
	var blocks []string
	var attachments []string

	// Will read the whole DATA mail and add it to a block called "mime", this is
	// only done to have a fully compatible email bridge
	fullMime, _ := ioutil.ReadAll(r)

	decodedMessage, err := common.DecodeFromMime(string(fullMime))
	if err != nil {
		return err
	}

	// Fetch both sender and recipient info
	svc := container.Instance.GetResolveService()

	// Check from address
	var fromAddr *address.Address
	if s.IsGateway {
		fromAddr, err = address.NewAddress(s.Account)
	} else {
		fromAddr, err = address.NewAddress(s.From)
	}
	if err != nil {
		return err
	}

	senderInfo, err := svc.ResolveAddress(fromAddr.Hash())
	if err != nil {
		return err
	}

	// Check to address
	toAddr, err := address.NewAddress(s.To)
	if err != nil {
		return err
	}

	recipientInfo, err := svc.ResolveAddress(toAddr.Hash())
	if err != nil {
		return err
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(fromAddr, nil, decodedMessage.From.Name, s.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	blocks = decodedMessage.Blocks
	if isGatewayDestination(s.To) {
		if decodedMessage.To == nil || len(decodedMessage.To) < 1 {
			return errIncorrectFormat
		}
		blocks = append(blocks, common.DestinationBlock+","+decodedMessage.To[0].Address)
	}

	for filename, base64Data := range decodedMessage.Attachments {
		b, _ := internal.Decode(base64Data)
		// We write the attachments temporary to disk so we can use it later on message.Compose,
		// however this needs to be improved so we don't need to write them to disk
		fName := filepath.Join(os.TempDir(), filename)
		err = ioutil.WriteFile(fName, b, 0644)
		if err != nil {
			return err
		}

		defer os.Remove(fName)

		attachments = append(attachments, fName)
	}

	// Compose mail
	envelope, err := message.Compose(addressing, decodedMessage.Subject, blocks, attachments)
	if err != nil {
		return err
	}

	// Send mail
	client, err := api.NewAuthenticated(*fromAddr, s.Info.GetActiveKey().PrivKey, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return err
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}

	return nil
}

// Reset is called when resetting the session
func (s *Session) Reset() {}

// Logout is called when logging out
func (s *Session) Logout() error {
	return nil
}

// isEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	if isEmailFromBitMaelumNetwork(e) {
		return true
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
