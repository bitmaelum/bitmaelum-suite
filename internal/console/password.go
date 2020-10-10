package console

import (
	"bytes"

	"github.com/chzyer/readline"
	"github.com/zalando/go-keyring"
)

const (
	service string = "bitmaelum"
	user    string = "vault"
)

var (
	// Readline instance. Override for mocking purposes
	readliner, _ = readline.NewEx(&readline.Config{
		DisableAutoSaveHistory: true,
	})
	// Override for mocking purposes
	kr keyring.Keyring
)

// StorePassword will store the given password into the keychain if possible
func StorePassword(pwd string) error {
	if kr != nil {
		return kr.Set(service, user, pwd)
	}

	return nil
}

// AskDoublePassword will ask for a password (and confirmation) on the commandline
func AskDoublePassword() ([]byte, error) {
	for {
		p1, _ := readliner.ReadPassword("Please enter your vault password: ")
		p2, _ := readliner.ReadPassword("Please retype your vault password: ")

		if bytes.Equal(p1, p2) {
			return p1, nil
		}

		_, _ = readliner.Stdout().Write([]byte("Passwords do not match. Please type again.\n"))
	}
}

// AskPassword will ask for a password (without confirmation) on the commandline
func AskPassword() (string, bool) {
	if kr != nil {
		pwd, err := kr.Get(service, user)
		if err == nil {
			return pwd, true
		}
	}

	b, _ := readliner.ReadPassword("Please enter your vault password: ")
	return string(b), false
}
