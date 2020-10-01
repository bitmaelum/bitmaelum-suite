package password

import (
	"bytes"
	"fmt"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

const (
	service string = "bitmaelum"
	user    string = "vault"
)

// StorePassword will store the given password into the keychain if possible
func StorePassword(pwd string) error {
	return keyring.Set(service, user, pwd)
}

// AskDoublePassword will ask for a password (and confirmation) on the commandline
func AskDoublePassword() []byte {
	for {
		fmt.Printf("Please enter your vault password: ")
		p1, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		fmt.Printf("\n")

		fmt.Printf("Please retype your vault password: ")
		p2, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		fmt.Printf("\n")

		if bytes.Equal(p1, p2) {
			return p1
		}

		fmt.Printf("Passwords do not match. Please type again.\n")
	}
}

// AskPassword will ask for a password (without confirmation) on the commandline
func AskPassword() (string, bool) {
	pwd, err := keyring.Get(service, user)
	if err == nil {
		return pwd, true
	}

	fmt.Printf("Please enter your vault password: ")
	b, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")

	return string(b), false
}
