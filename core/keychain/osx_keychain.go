package keychain

import (
    "errors"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/keybase/go-keychain"
)

const (
    ACCESSGROUP string = "keychain.bitmaelum.nl"
    SERVICE     string = "bitmaelum"
)

type OSXKeyChain struct {
}

var (
	// Key not found
	ErrorKeyNotFound = errors.New("key not found")
)

func (kc *OSXKeyChain) Fetch(addr core.Address) ([]byte, error) {
    query := keychain.NewItem()
    query.SetSecClass(keychain.SecClassGenericPassword)
    query.SetService(SERVICE)
    query.SetAccount(addr.String())
    query.SetAccessGroup(ACCESSGROUP)
    query.SetMatchLimit(keychain.MatchLimitOne)
    query.SetReturnData(true)
    results, err := keychain.QueryItem(query)
    if err != nil {
      return nil, err
    } else if len(results) != 1 {
      return nil, ErrorKeyNotFound
    } else {
      return results[0].Data, nil
    }
}

func (kc *OSXKeyChain) Store(addr core.Address, key []byte) error {
    item := keychain.NewItem()
    item.SetSecClass(keychain.SecClassGenericPassword)
    item.SetService(SERVICE)
    item.SetAccount(addr.String())
    item.SetLabel("")
    item.SetAccessGroup(ACCESSGROUP)
    item.SetData(key)
    item.SetSynchronizable(keychain.SynchronizableNo)
    item.SetAccessible(keychain.AccessibleWhenUnlocked)

    return keychain.AddItem(item)
}
