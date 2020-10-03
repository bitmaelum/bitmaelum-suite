package address

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

var (
	hashRegex    = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

// Hash is a SHA256'd address
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
