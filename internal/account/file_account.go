package account

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
)

// // CreateAccount creates new account for the given address and public key
// func (r *fileRepo) Create(addr address.HashAddress, pubKey string) error {
// 	if r.Exists(addr) {
// 		return errors.New("account already exists")
// 	}
//
// 	_ = r.CreateBox(addr, internal.BoxInbox)
// 	_ = r.CreateBox(addr, internal.BoxOutbox)
// 	_ = r.CreateBox(addr, internal.BoxTrash)
// 	_ = r.StorePubKey(addr, pubKey)
//
// 	return nil
// }

// Create a new account for this address
func (r *fileRepo) Create(addr address.HashAddress, pubKeyPEM string) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	err := os.MkdirAll(fullPath, 0700)
	if err != nil {
		return err
	}

	// parallelize actions
	g := new(errgroup.Group)
	g.Go(func() error {
		return r.StoreKey(addr, pubKeyPEM)
	})
	for _, box := range MandatoryBoxes {
		g.Go(func() error {
			return r.CreateBox(addr, box)
		})
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		_ = r.Delete(addr)
		return err
	}

	return nil
}

// Returns true when the given account for this address exists
func (r *fileRepo) Exists(addr address.HashAddress) bool {
	return r.pathExists(addr, "")
}

// Delete an account
func (r *fileRepo) Delete(addr address.HashAddress) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	return os.RemoveAll(fullPath)
}
