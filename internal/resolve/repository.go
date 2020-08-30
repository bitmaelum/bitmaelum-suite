package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

var errKeyNotFound = errors.New("hash not found")

// Repository is the interface to manage address resolving
type Repository interface {
	Resolve(addr address.HashAddress) (*Info, error)
	Upload(addr address.HashAddress, pubKey encrypt.PubKey, resolveAddress, signature string) error
}
