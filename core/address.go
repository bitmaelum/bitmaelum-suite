package core

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "regexp"
    "strings"
)

const (
    // This is the main regex where an address should confirm to. Much simpler than an email address
    ADDRESS_REGEX   string = "(^[a-z0-9][a-z0-9\\.\\-]{2,63})(?:@([a-z0-9][a-z0-9\\.\\-]{1,63}))?!$"
)


// SHA256'd address
type HashAddress    string

// An actual address
type Address struct {
    Local  string
    Org    string
}

// Actual hash function
func StringToHash(address string) string {
    sum := sha256.Sum256([]byte(address))
    return hex.EncodeToString(sum[:])
}

// Return true when the given string is a valid address
func IsValidAddress(address string) bool {
    _, err := NewAddressFromString(address)
    return err == nil
}

// Returns a valid address structure
func NewAddressFromString(address string) (*Address, error) {
    re := regexp.MustCompile(ADDRESS_REGEX)
    if re == nil {
        return nil, errors.New("cannot compile regex")
    }

    if ! re.MatchString(strings.ToLower(address)) {
        return nil, errors.New("incorrect format specified")
    }

    matches := re.FindStringSubmatch(strings.ToLower(address))

    return &Address{
        Local: matches[1],
        Org: matches[2],
    }, nil
}

// Converts an address to a string
func (a *Address) ToString() string {
    if a.Org == "" {
        return fmt.Sprintf("%s!", a.Local)
    }

    return fmt.Sprintf("%s@%s!", a.Local, a.Org)
}

// Converts an address to a hashed value
func (a *Address) ToHash() string {
    return StringToHash(a.ToString())
}
