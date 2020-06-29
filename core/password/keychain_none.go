// +build windows linux

package password

import "github.com/bitmaelum/bitmaelum-server/core"

func (kc *KeyChain) IsAvailable() bool {
	return false
}

func (kc *KeyChain) Fetch(addr core.Address) ([]byte, error) {
	return nil, nil
}

func (kc *KeyChain) Store(addr core.Address, key []byte) error {
	return nil
}
