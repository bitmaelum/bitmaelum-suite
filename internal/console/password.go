package console

import (
	"bytes"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	service string = "bitmaelum"
	user    string = "vault"
)

var (
	// Override for mocking purposes
	kr keyring.Keyring
	pwdReader PasswordReader
)

// PasswordReader is an interface to read a password.
type PasswordReader interface {
    ReadPassword() ([]byte, error)
}

// StdInPasswordReader Default password reader from stdin
type StdInPasswordReader struct {
}

// ReadPassword reads password from stdin
func (pr *StdInPasswordReader) ReadPassword() ([]byte, error) {
    pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
    return pwd, err
}


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
		fmt.Print("Please enter your vault password: ")
		p1, err := pwdReader.ReadPassword()
		if err != nil {
			return nil, err
		}

		fmt.Print("Please retype your vault password: ")
		p2, err := pwdReader.ReadPassword()
		if err != nil {
			return nil, err
		}

		if bytes.Equal(p1, p2) {
			return p1, nil
		}

		fmt.Println("Passwords do not match. Please type again.")
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

	fmt.Print("Please enter your vault password: ")
	b, _ := pwdReader.ReadPassword()
	return string(b), false
}
