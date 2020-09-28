package address

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// This is the main regex where an address should confirm to. Much simpler than an email address
	addressRegex = regexp.MustCompile("(^[a-z0-9][a-z0-9\\.\\-]{2,63})(?:@([a-z0-9][a-z0-9\\.\\-]{1,63}))?!$")
	orgRegex     = regexp.MustCompile("^[a-z0-9][a-z0-9\\.\\-]{1,63}$")
	hashRegex    = regexp.MustCompile("[a-f0-9]{64}")
)

// HashAddress is a SHA256'd address
type HashAddress string

// HashOrganisation is a SHA256'd organisation
type HashOrganisation string

// String casts an hash address to string
func (ha HashAddress) String() string {
	return string(ha)
}

// String casts an hash address to string
func (ha HashOrganisation) String() string {
	return string(ha)
}

// Address represents a bitMaelum address
type Address struct {
	Local string
	Org   string
}

// New returns a valid address structure based on the given address
func New(address string) (*Address, error) {
	if !addressRegex.MatchString(strings.ToLower(address)) {
		return nil, errors.New("incorrect address format specified")
	}

	matches := addressRegex.FindStringSubmatch(strings.ToLower(address))

	return &Address{
		Local: matches[1],
		Org:   matches[2],
	}, nil
}

// NewHash generates a hash address based on the given email string
func NewHash(address string) (*HashAddress, error) {
	a, err := New(address)
	if err != nil {
		return nil, err
	}

	h := a.Hash()
	return &h, nil
}

// NewHashFromHash generates a hash address based on the given string hash
func NewHashFromHash(hash string) (*HashAddress, error) {
	if !hashRegex.MatchString(strings.ToLower(hash)) {
		return nil, errors.New("incorrect hash address format specified")
	}

	h := HashAddress(hash)
	return &h, nil
}

// IsValidAddress returns true when the given string is a valid BitMaelum address
func IsValidAddress(address string) bool {
	_, err := New(address)
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
func (a *Address) Hash() HashAddress {
	l := sha256.Sum256([]byte(a.Local))
	o := sha256.Sum256([]byte(a.Org))
	sum := sha256.Sum256([]byte(hex.EncodeToString(l[:]) + hex.EncodeToString(o[:])))

	return HashAddress(hex.EncodeToString(sum[:]))
}

// OrganisationHash converts an address to a hashed value
func (a *Address) OrganisationHash() HashOrganisation {
	l := sha256.Sum256([]byte(""))
	o := sha256.Sum256([]byte(a.Org))
	sum := sha256.Sum256([]byte(hex.EncodeToString(l[:]) + hex.EncodeToString(o[:])))

	return HashOrganisation(hex.EncodeToString(sum[:]))
}

// NewOrgHash returns a new organisation hash
func NewOrgHash(org string) (*HashOrganisation, error) {
	if !orgRegex.MatchString(strings.ToLower(org)) {
		return nil, errors.New("incorrect org format specified")
	}

	a := &Address{
		Local: "",
		Org:   strings.ToLower(org),
	}

	h := a.OrganisationHash()
	return &h, nil
}

// Bytes converts an address to []byte
func (a *Address) Bytes() []byte {
	return []byte(a.String())
}

// VerifyHash will check if the hashes for local and org found matches the actual target hash
func VerifyHash(target, local, org string) bool {
	sum := sha256.Sum256([]byte(local + org))
	hash := hex.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(hash), []byte(target)) == 1
}
