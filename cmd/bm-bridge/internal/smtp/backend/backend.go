// Copyright (c) 2022 BitMaelum Authors
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
	"strings"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-smtp"
	"github.com/sirupsen/logrus"
)

var (
	errIncorrectFormat = errors.New("incorrect account format specified")
	errAccountNotFound = errors.New("account not found")
	errGatewayOnly     = errors.New("this is a email->bitmaelum gateway only")
)

// Backend will hold the Vault to use and the GatewayAccount if needed
type Backend struct {
	Vault          *vault.Vault
	GatewayAccount string
}

// New will create a new backend
func New(v *vault.Vault, ga string) *Backend {
	return &Backend{
		Vault:          v,
		GatewayAccount: ga,
	}
}

// Login handles a login command with username and password. We accept any password as long the
// username is a valid account on the vault
func (be *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if be.GatewayAccount != "" {
		return nil, errGatewayOnly
	}

	if !strings.HasSuffix(username, "!") {
		username = username + "!"
	}

	session, err := be.getSessionFromAccount(username)

	if err != nil {
		logrus.Errorf("SMTP: user %s error when login - %s", username, err.Error())
	} else {
		logrus.Infof("SMTP: user %s logged in", username)
	}

	return session, err
}

// AnonymousLogin will accept mails if this server allows to be a email<->bitmaelum gateway
func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	if be.GatewayAccount != "" {
		session, err := be.getSessionFromAccount(be.GatewayAccount)
		if err != nil {
			return nil, err
		}

		session.IsGateway = true
		session.RemoteAddr = state.RemoteAddr

		logrus.Infof("SMTP: external connection received from %s", state.RemoteAddr.String())

		return session, err
	}

	return nil, smtp.ErrAuthRequired
}

func (be *Backend) getSessionFromAccount(acc string) (*Session, error) {
	session := &Session{}

	addr, err := address.NewAddress(acc)
	if err != nil {
		return session, errIncorrectFormat
	}

	if !be.Vault.HasAccount(*addr) {
		return session, errAccountNotFound
	}

	// If account exists in vault, set up a new session
	session.Account = acc
	session.Vault = be.Vault
	session.Info, session.Client, err = common.GetClientAndInfo(be.Vault, acc)
	if err != nil {
		return session, err
	}

	return session, nil
}
