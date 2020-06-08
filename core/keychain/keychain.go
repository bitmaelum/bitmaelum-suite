package keychain

import "github.com/jaytaph/mailv2/core"

type KeyChain interface {
    // Fetch a key from the keychain
    Fetch(addr core.Address) ([]byte, error)

    // Fetch a key form the keychain
    Store(addr core.Address, key []byte) error
}
