package resolve

import (
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
)

// Repository is the interface to manage address resolving
type Repository interface {
	Resolve(addr address.HashAddress) (*Info, error)
	Upload(addr address.HashAddress, pubKey, resolveAddress, signature string) error
}
