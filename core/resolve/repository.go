package resolve

import "github.com/bitmaelum/bitmaelum-server/core"

// Repository is the interface to manage address resolving
type Repository interface {
	Resolve(addr core.HashAddress) (*Info, error)
	Upload(addr core.HashAddress, pubKey, resolveAddress, signature string) error
}
