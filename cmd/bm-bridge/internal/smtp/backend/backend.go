package smtpgw

import (
	"errors"

	common "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/emersion/go-smtp"
)

type Backend struct {
	Vault *vault.Vault
}

func New(v *vault.Vault) *Backend {
	return &Backend{
		Vault: v,
	}
}

// Login handles a login command with username and password.
func (be *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	account := username + "!"

	addr, err := address.NewAddress(account)
	if err != nil {
		return nil, errors.New("NO incorrect address format specified")
	}
	if !be.Vault.HasAccount(*addr) {
		return nil, errors.New("NO account not found in vault")
	}

	session := &Session{Account: account, Vault: be.Vault}

	session.Info, session.Client, err = common.GetClientAndInfo(be.Vault, account)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	session := &Session{Account: "", Vault: be.Vault}

	return session, nil
	//return nil, smtp.ErrAuthRequired
}
