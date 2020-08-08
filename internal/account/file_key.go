package account

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/nightlyone/lockfile"
)

// Store the public key for this account
func (r *fileRepo) StoreKey(addr address.HashAddress, key string) error {
	// Lock our key file for writing
	lockfilePath := r.getPath(addr, pubKeyFile+".lock")
	lock, err := lockfile.New(lockfilePath)
	if err != nil {
		return err
	}

	err = lock.TryLock()
	if err != nil {
		return err
	}

	defer func() {
		_ = lock.Unlock()
	}()

	// Read keys
	pk := &PubKeys{}
	err = r.fetchJSON(addr, pubKeyFile, pk)
	if err != nil {
		return err
	}

	// Add new key
	pk.PubKeys = append(pk.PubKeys, key)

	// Convert back to string
	data, err := json.MarshalIndent(pk, "", "  ")
	if err != nil {
		return err
	}

	// And store
	return r.store(addr, pubKeyFile, data)
}

// Retrieve the public keys for this account
func (r *fileRepo) FetchKeys(addr address.HashAddress) ([]string, error) {
	pk := &PubKeys{}
	err := r.fetchJSON(addr, pubKeyFile, pk)
	if err != nil {
		return nil, err
	}

	return pk.PubKeys, nil
}

// Retrieve the public keys for this account and decode the PEMs
func (r *fileRepo) FetchDecodedKeys(addr address.HashAddress) ([]interface{}, error) {
	keys, err := r.FetchKeys(addr)
	if err != nil {
		return nil, err
	}

	s := make([]interface{}, len(keys))
	for i, pem := range keys {
		k, err := encrypt.PEMToPubKey([]byte(pem))
		if err != nil {
			continue
		}

		s[i] = k
	}

	return s, nil
}
