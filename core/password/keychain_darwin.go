// +build darwin

package password

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-server/core"
	gokeychain "github.com/keybase/go-keychain"
)

const (
	accessGroup string = "keychain.bitmaelum.nl"
	service     string = "bitmaelum"
)

var (
	// ErrorKeyNotFound is returned when a key could not be found in the keychain
	ErrorKeyNotFound = errors.New("key not found")
)

// IsAvailable returns true when a keychain is available
func (kc *KeyChain) IsAvailable() bool {
	return true
}

// Fetch returns a key/password for the given address
func (kc *KeyChain) Fetch(addr core.Address) ([]byte, error) {
	query := gokeychain.NewItem()
	query.SetSecClass(gokeychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(addr.String())
	query.SetAccessGroup(accessGroup)
	query.SetMatchLimit(gokeychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := gokeychain.QueryItem(query)
	if err != nil {
		return nil, err
	} else if len(results) != 1 {
		return nil, ErrorKeyNotFound
	} else {
		return results[0].Data, nil
	}
}

// Store stores a key/password for the given address in the keychain/vault
func (kc *KeyChain) Store(addr core.Address, key []byte) error {
	item := gokeychain.NewItem()
	item.SetSecClass(gokeychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(addr.String())
	item.SetLabel("")
	item.SetAccessGroup(accessGroup)
	item.SetData(key)
	item.SetSynchronizable(gokeychain.SynchronizableNo)
	item.SetAccessible(gokeychain.AccessibleAfterFirstUnlockThisDeviceOnly)

	return gokeychain.AddItem(item)
}
