package account

import (
	"encoding/json"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
)

// Store the public key for this account
func (r *fileRepo) StoreKey(addr address.Hash, key bmcrypto.PubKey) error {
	// Lock our key file for writing
	lockfilePath := r.getPath(addr, keysFile+".lock")
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
	err = r.fetchJSON(addr, keysFile, pk)
	if err != nil && os.IsNotExist(err) {
		err = r.createPubKeyFile(addr)
	}
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
	return r.store(addr, keysFile, data)
}

// Retrieve the public keys for this account
func (r *fileRepo) FetchKeys(addr address.Hash) ([]bmcrypto.PubKey, error) {
	pk := &PubKeys{}
	err := r.fetchJSON(addr, keysFile, pk)
	if err != nil && os.IsNotExist(err) {
		err = r.createPubKeyFile(addr)
	}
	if err != nil {
		return nil, err
	}

	return pk.PubKeys, nil
}

// Create file because it doesn't exist yet
func (r *fileRepo) createPubKeyFile(addr address.Hash) error {
	p := r.getPath(addr, keysFile)
	f, err := os.Create(p)
	if err != nil {
		logrus.Errorf("Error while creating file %s: %s", p, err)
		return err
	}
	return f.Close()
}
