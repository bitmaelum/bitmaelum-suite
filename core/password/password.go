package password

import (
	"bytes"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

// FetchPassword tries to figure out the password of the given address. It can do so by checking
// keychains and such, or when all fails, ask it from the user.
func FetchPassword(addr *address.Address) ([]byte, error) {
	if keychain.IsAvailable() {
		pwd, err := keychain.Fetch(*addr)
		if err == nil {
			return pwd, nil
		}
	}

	// If all fails, ask from stdin
	fmt.Printf("\U0001F511  Please enter your password for account '%s': ", addr.String())
	p, e := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")

	return p, e
}

// StorePassword will store the given password into the keychain if possible
func StorePassword(addr *address.Address, pwd []byte) error {
	if keychain.IsAvailable() {
		_, err := keychain.Fetch(*addr)
		if err != nil {
			return keychain.Store(*addr, pwd)
		}
	}

	return nil
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

		if bytes.Compare(p1, p2) == 0 {
			return p1
		}

		fmt.Printf("Passwords do not match. Please type again.\n")
	}
}

// AskPassword will ask for a password (without confirmation) on the commandline
func AskPassword() []byte {
	fmt.Printf("Please enter your vault password: ")
	p1, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")

	return p1
}
