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

package console

import (
	"bytes"
	"fmt"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	service string = "bitmaelum"
	user    string = "vault"
)

var (
	// Override for mocking purposes
	kr        keyring.Keyring
	pwdReader PasswordReader = &StdInPasswordReader{}
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
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
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
func AskDoublePassword() (string, error) {
	return AskDoublePasswordPrompt("Please enter your new vault password: ", "Please retype your new vault password: ")
}

// AskDoublePasswordPrompt will ask for a password (and confirmation) on the commandline, with a custom prompt
func AskDoublePasswordPrompt(p1, p2 string) (string, error) {
	for {
		fmt.Print(p1)
		p1, err := pwdReader.ReadPassword()
		if err != nil {
			return "", err
		}
		fmt.Println("")

		fmt.Print(p2)
		p2, err := pwdReader.ReadPassword()
		if err != nil {
			return "", err
		}
		fmt.Println("")

		if bytes.Equal(p1, p2) {
			return string(p1), nil
		}

		fmt.Println("Passwords do not match. Please try again.")
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
	fmt.Println("")
	fmt.Println("")

	return string(b), false
}
