package address

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// This is the main regex where an address should confirm to. Much simpler than an email address
	addressRegex = regexp.MustCompile(`^([a-z0-9][a-z0-9.\-]{1,62}[a-z0-9])(?:@([a-z0-9][a-z0-9.\-]{0,62}[a-z0-9]))?!$`)
)

// Address represents a bitMaelum address
type Address struct {
	Local string
	Org   string
}

// New returns a valid address structure based on the given address
func NewAddress(address string) (*Address, error) {
	if !addressRegex.MatchString(strings.ToLower(address)) {
		return nil, errors.New("incorrect address format specified")
	}

	matches := addressRegex.FindStringSubmatch(strings.ToLower(address))

	return &Address{
		Local: matches[1],
		Org:   matches[2],
	}, nil
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
func (a *Address) Hash() Hash {
	return NewHash(a.LocalHash().String() + a.OrgHash().String())
}

// LocalHash returns the hash of the local/user part of the address
func (a *Address) LocalHash() Hash {
	return NewHash(a.Local)
}

// OrgHash returns the hash of the organisation part of the address
func (a *Address) OrgHash() Hash {
	return NewHash(a.Org)
}

// Bytes converts an address to []byte
func (a *Address) Bytes() []byte {
	return []byte(a.String())
}
