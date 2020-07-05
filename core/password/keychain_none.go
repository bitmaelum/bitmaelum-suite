// +build windows linux

package password

import (
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
)

// IsAvailable returns true when a keychain is available
func (kc *KeyChain) IsAvailable() bool {
	return false
}

// Fetch returns a key/password for the given address
func (kc *KeyChain) Fetch(addr address.Address) ([]byte, error) {
	return nil, nil
}

// Store stores a key/password for the given address in the keychain/vault
func (kc *KeyChain) Store(addr address.Address, key []byte) error {
	return nil
}
