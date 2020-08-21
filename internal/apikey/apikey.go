package apikey

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/xhit/go-str2duration/v2"
	"strconv"
	"strings"
	"time"
)

// KeyType represents a key with a validity and permissions. When admin is true, permission checks are always true
type KeyType struct {
	ID          string    `json:"key"`
	ValidUntil  time.Time `json:"valid_until"`
	Permissions []string  `json:"permissions"`
	Admin       bool      `json:"admin"`
}

// NewAdminKey creates a new admin key
func NewAdminKey(valid time.Duration) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          internal.GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: nil,
		Admin:       true,
	}
}

// NewKey creates a new key with the given permissions and duration
func NewKey(perms []string, valid time.Duration) KeyType {
	var until = time.Time{}
	if valid > 0 {
		until = time.Now().Add(valid)
	}

	return KeyType{
		ID:          internal.GenerateKey("BMK-", 32),
		ValidUntil:  until,
		Permissions: perms,
		Admin:       false,
	}
}

// HasPermission returns true when this key contains the given permission, or when the key is an admin key (granted all permissions)
func (key *KeyType) HasPermission(perm string) bool {
	if key.Admin == true {
		return true
	}

	perm = strings.ToLower(perm)

	for _, p := range key.Permissions {
		if strings.ToLower(p) == perm {
			return true
		}
	}

	return false
}

// ParsePermissions checks all permission and returns an error when a permission is not valid
func ParsePermissions(perms []string) error {
	for _, p := range perms {
		p = strings.ToLower(p)

		found := false
		for _, ap := range AllPermissons {
			if p == ap {
				found = true
			}
		}
		if !found {
			return errors.New("unknown permission: " + p)
		}
	}

	return nil
}

// ParseValidDuration gets a time duration string and return the time duration. Accepts single int as days
func ParseValidDuration(ds string) (time.Duration, error) {
	vd, err := str2duration.ParseDuration(ds)
	if err != nil {
		days, err := strconv.Atoi(ds)
		if err != nil {
			return 0, err
		}

		vd = time.Duration(days*24) * time.Hour
	}

	return vd, nil
}

// Repository is a repository to fetch and store API keys
type Repository interface {
	Fetch(ID string) (*KeyType, error)
	Store(key KeyType) error
	Remove(ID string)
}
