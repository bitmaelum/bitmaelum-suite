package apikey

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// KeyType represents a key with a validity and permissions. When admin is true, permission checks are always true
type KeyType struct {
	ID          string     `json:"key"`
	ValidUntil  time.Time  `json:"valid_until"`
	Permissions []string   `json:"permissions"`
	Admin       bool       `json:"admin"`
	AddrHash    *hash.Hash `json:"addr_hash"`
	Desc        string     `json:"description"`
}

// NewAdminKey creates a new admin key
func NewAdminKey(valid time.Duration, desc string) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: nil,
		Admin:       true,
		AddrHash:    nil,
		Desc:        desc,
	}
}

// NewAccountKey creates a new key with the given permissions and duration for the given account
func NewAccountKey(addrHash *hash.Hash, perms []string, valid time.Duration, desc string) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: perms,
		Admin:       false,
		AddrHash:    addrHash,
		Desc:        desc,
	}
}

// NewKey creates a new key with the given permissions and duration
func NewKey(perms []string, valid time.Duration, desc string) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: perms,
		Admin:       false,
		AddrHash:    nil,
		Desc:        desc,
	}
}

// HasPermission returns true when this key contains the given permission, or when the key is an admin key (granted all permissions)
func (key *KeyType) HasPermission(perm string, addrHash *hash.Hash) bool {
	if key.Admin {
		return true
	}

	perm = strings.ToLower(perm)

	for _, p := range key.Permissions {
		if strings.ToLower(p) == perm && (addrHash == nil || addrHash == key.AddrHash) {
			return true
		}
	}

	return false
}

// Repository is a repository to fetch and store API keys
type Repository interface {
	Fetch(ID string) (*KeyType, error)
	FetchByHash(h string) ([]*KeyType, error)
	Store(key KeyType) error
	Remove(key KeyType) error
}

// GenerateKey generates a random key based on a given string length
func GenerateKey(prefix string, n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return prefix + string(b)
}
