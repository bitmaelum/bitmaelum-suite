// +build darwin

package password

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-server/core"
	gokeychain "github.com/keybase/go-keychain"
)

const (
	ACCESSGROUP string = "keychain.bitmaelum.nl"
	SERVICE     string = "bitmaelum"
)

var (
	// Key not found
	ErrorKeyNotFound = errors.New("key not found")
)

func (kc *KeyChain) IsAvailable() bool {
	return true
}

func (kc *KeyChain) Fetch(addr core.Address) ([]byte, error) {
	query := gokeychain.NewItem()
	query.SetSecClass(gokeychain.SecClassGenericPassword)
	query.SetService(SERVICE)
	query.SetAccount(addr.String())
	query.SetAccessGroup(ACCESSGROUP)
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

func (kc *KeyChain) Store(addr core.Address, key []byte) error {
	item := gokeychain.NewItem()
	item.SetSecClass(gokeychain.SecClassGenericPassword)
	item.SetService(SERVICE)
	item.SetAccount(addr.String())
	item.SetLabel("")
	item.SetAccessGroup(ACCESSGROUP)
	item.SetData(key)
	item.SetSynchronizable(gokeychain.SynchronizableNo)
	item.SetAccessible(gokeychain.AccessibleAfterFirstUnlockThisDeviceOnly)

	return gokeychain.AddItem(item)
}
