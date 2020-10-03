package address

import (
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

// OrganisationHash is a SHA256'd organisation
type OrganisationHash string

// String casts an hash address to string
func (ha OrganisationHash) String() string {
	return string(ha)
}

// NewOrganisationHash returns a new organisation hash. This is the hash of the organisation part
func NewOrganisationHash(org string) (*OrganisationHash, error) {
	if !orgRegex.MatchString(strings.ToLower(org)) {
		return nil, errors.New("incorrect org format specified")
	}

	a := &Address{
		Local: "",
		Org:   strings.ToLower(org),
	}

	h := OrganisationHash(a.OrgHash())
	return &h, nil
}

// OrganisationHash returns the organisational hash generated from the org
func (a *Address) OrganisationHash() OrganisationHash {
	l := sha256.Sum256([]byte(""))
	o := sha256.Sum256([]byte(a.Org))
	sum := sha256.Sum256([]byte(hex.EncodeToString(l[:]) + hex.EncodeToString(o[:])))

	return OrganisationHash(hex.EncodeToString(sum[:]))
}

// HasOrganisationPart returns true when the address is an organisational address (user@org!)
func (a *Address) HasOrganisationPart() bool {
	return len(a.Org) > 0
}

// VerifyHash will check if the hashes for local and org found matches the actual target hash
func VerifyHash(targetHash, localHash, orgHash Hash) bool {
	h := NewHash(localHash.String() + orgHash.String())

	return subtle.ConstantTimeCompare(h.Byte(), targetHash.Byte()) == 1
}
