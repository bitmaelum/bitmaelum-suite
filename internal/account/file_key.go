package account

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"os"
)


todo: we must convert from pubkey

// Store the public key for this account
func (r *fileRepo) StoreKey(addr address.HashAddress, key encrypt.PubKey) error {
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
	return r.store(addr, pubKeyFile, data)
}



// Retrieve the public keys for this account
func (r *fileRepo) FetchKeys(addr address.HashAddress) ([]encrypt.PubKey, error) {
	pk := &PubKeys{}
	err := r.fetchJSON(addr, pubKeyFile, pk)
	if err != nil && os.IsNotExist(err) {
		err = r.createPubKeyFile(addr)
	}
	if err != nil {
		return nil, err
	}

	return pk.PubKeys, nil
}

// Create file because it doesn't exist yet
func (r *fileRepo) createPubKeyFile(addr address.HashAddress) error {
	p := r.getPath(addr, pubKeyFile)
	f, err := os.Create(p)
	if err != nil {
		logrus.Errorf("Error while creating file %s: %s", p, err)
		return err
	}
	return f.Close()
}
