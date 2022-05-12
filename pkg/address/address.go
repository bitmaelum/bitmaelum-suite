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

package address

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

var errIncorrectAddress = func(addr string) error { return fmt.Errorf("incorrect address format specified: %s", addr) }

var (
	// This is the main regex where an address should confirm to. Much simpler than an email address
	addressRegex = regexp.MustCompile(`^([a-z0-9][a-z0-9.\-]{1,62}[a-z0-9])(?:@([a-z0-9][a-z0-9.\-]{0,62}[a-z0-9]))?!$`)
)

// Address represents a bitMaelum address
type Address struct {
	Local string // Local part is either <local>!  or  <local>@<organisation>!
	Org   string // Org part is either "" in case of <local>!  or <local>@<organisation>!
}

// UnmarshalJSON will unmarshal a string of bytes into an address
func (a *Address) UnmarshalJSON(b []byte) error {
	var s string

	// Decode string first
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	tmp, err := NewAddress(s)
	if err != nil {
		return err
	}

	// Copy data
	a.Local = tmp.Local
	a.Org = tmp.Org

	return nil
}

// MarshalJSON will marshal an address into a string of bytes for JSON output
func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// NewAddress returns a valid address structure based on the given address
func NewAddress(address string) (*Address, error) {
	address = sanitizeAddress(address)

	if !addressRegex.MatchString(address) {
		return nil, errIncorrectAddress(address)
	}

	matches := addressRegex.FindStringSubmatch(strings.ToLower(address))

	return &Address{
		Local: matches[1],
		Org:   matches[2],
	}, nil
}

// Sanitizes the address by removing any dashes and dots
func sanitizeAddress(address string) string {
	address = strings.ToLower(address)
	address = strings.ReplaceAll(address, ".", "")
	address = strings.ReplaceAll(address, "-", "")
	return address
}

// IsValidAddress returns true when the given string is a valid BitMaelum address
func IsValidAddress(address string) bool {
	_, err := NewAddress(address)
	return err == nil
}

// String converts an address to a string
func (a *Address) String() string {
	if len(a.Org) == 0 {
		return fmt.Sprintf("%s!", a.Local)
	}

	return fmt.Sprintf("%s@%s!", a.Local, a.Org)
}

// Hash converts an address to a hashed value
func (a *Address) Hash() hash.Hash {
	return hash.New(a.LocalHash().String() + a.OrgHash().String())
}

// LocalHash returns the hash of the local/user part of the address
func (a *Address) LocalHash() hash.Hash {
	return hash.New(a.Local)
}

// OrgHash returns the hash of the organisation part of the address
func (a *Address) OrgHash() hash.Hash {
	return hash.New(a.Org)
}

// HasOrganisationPart returns true when the address is an organisational address (user@org!)
func (a *Address) HasOrganisationPart() bool {
	return len(a.Org) > 0
}

// Bytes converts an address to []byte
func (a *Address) Bytes() []byte {
	return []byte(a.String())
}
