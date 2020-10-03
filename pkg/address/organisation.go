package address

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var (
	orgRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.\-]{0,62}[a-z0-9]$`)
)

// HashOrganisation is a SHA256'd organisation
type HashOrganisation string

// String casts an hash address to string
func (ha HashOrganisation) String() string {
	return string(ha)
}

// NewOrgHash returns a new organisation hash. This is the hash of the organisation part
func NewOrgHash(org string) (*HashOrganisation, error) {
	if !orgRegex.MatchString(strings.ToLower(org)) {
		return nil, errors.New("incorrect org format specified")
	}

	a := &Address{
		Local: "",
		Org:   strings.ToLower(org),
	}

	h := HashOrganisation(a.OrgHash())
	return &h, nil
}

// OrganisationHash converts an address to a hashed value
func (a *Address) OrganisationHash() HashOrganisation {
	l := sha256.Sum256([]byte(""))
	o := sha256.Sum256([]byte(a.Org))
	sum := sha256.Sum256([]byte(hex.EncodeToString(l[:]) + hex.EncodeToString(o[:])))

	return HashOrganisation(hex.EncodeToString(sum[:]))
}

// HasOrganisationPart returns true when the address is an organisational address (user@org!)
func (a *Address) HasOrganisationPart() bool {
	return len(a.Org) > 0
}

// VerifyHash will check if the hashes for local and org found matches the actual target hash
func VerifyHash(targetHash, localHash, orgHash Hash) bool {
	l, err := hex.DecodeString(localHash.String())
	if err != nil {
		return false
	}
	o, err := hex.DecodeString(orgHash.String())
	if err != nil {
		return false
	}

	sum := sha256.Sum256(bytes.Join([][]byte{l, o}, []byte{}))
	h := NewHash(hex.EncodeToString(sum[:]))

	return subtle.ConstantTimeCompare(h.Byte(), targetHash.Byte()) == 1
}
