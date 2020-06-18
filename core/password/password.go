package password

import (
	"bytes"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"golang.org/x/crypto/ssh/terminal"
	"runtime"
	"syscall"
)

// FetchPassword tries to figure out the password of the given address. It can do so by checking
// keychains and such, or when all fails, ask it from the user.
func FetchPassword(addr *core.Address) ([]byte, error) {
	// Check OSX keychain if available
	if runtime.GOOS == "darwin" {
		keychain := OSXKeyChain{}
		pwd, err := keychain.Fetch(*addr)
		if err == nil {
			return pwd, nil
		}
	}

	// If all fails, ask from stdin
	fmt.Printf("\U0001F511  Please enter your password for account '%s': ", addr.String())
	p, e := terminal.ReadPassword(syscall.Stdin)
	fmt.Printf("\n")

	return p, e
}

func StorePassword(addr *core.Address, pwd []byte) error {
	// Check OSX keychain if available
	if runtime.GOOS == "darwin" {
		keychain := OSXKeyChain{}
		_, err := keychain.Fetch(*addr)
		if err != nil {
			keychain := OSXKeyChain{}
			return keychain.Store(*addr, pwd)
		}
	}

	return nil
}

func AskPassword() []byte {
	for {
		fmt.Printf("Please enter your password: ")
		p1, _ := terminal.ReadPassword(syscall.Stdin)
		fmt.Printf("\n")

		fmt.Printf("Please retype your password: ")
		p2, _ := terminal.ReadPassword(syscall.Stdin)
		fmt.Printf("\n")

		if bytes.Compare(p1, p2) == 0 {
			return p1
		}

		fmt.Printf("Passwords do not match. Please type again.\n")
	}
}