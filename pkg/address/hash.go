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
	hashRegex = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

// Hash is a hashed entity. This could be an address, localpart, organisation part or anything else. Currently only
// sha256 is supported
type Hash string

// NewHash generates a regular hash. Assumes you know what you are hashing
func NewHash(s string) Hash {
	sum := sha256.Sum256([]byte(s))

	return Hash(hex.EncodeToString(sum[:]))
}

// HashFromString generates a hash address based on the given string hash
func HashFromString(hash string) (*Hash, error) {
	if !hashRegex.MatchString(strings.ToLower(hash)) {
		return nil, errors.New("incorrect hash address format specified")
	}

	h := Hash(hash)
	return &h, nil
}

// String casts an hash address to string
func (ha Hash) String() string {
	return string(ha)
}

// Byte casts an hash address to string
func (ha Hash) Byte() []byte {
	return []byte(ha)
}

// VerifyHash will check if the hashes for local and org found matches the actual target hash
func (ha Hash) Verify(localHash, orgHash Hash) bool {
	target := NewHash(localHash.String() + orgHash.String())

	return subtle.ConstantTimeCompare(ha.Byte(), target.Byte()) == 1
}
