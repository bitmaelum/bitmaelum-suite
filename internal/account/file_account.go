package account

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Create a new account for this address
func (r *fileRepo) Create(addr hash.Hash, pubKey bmcrypto.PubKey) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	err := os.MkdirAll(fullPath, 0700)
	if err != nil {
		return err
	}

	// parallelize actions
	g := new(errgroup.Group)
	g.Go(func() error {
		return r.StoreKey(addr, pubKey)
	})
	for _, box := range MandatoryBoxes {
		boxCopy := box
		g.Go(func() error {
			logrus.Trace("Creating box: ", boxCopy)
			return r.CreateBox(addr, boxCopy)
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
func (r *fileRepo) Exists(addr hash.Hash) bool {
	return r.pathExists(addr, "")
}

// Delete an account
func (r *fileRepo) Delete(addr hash.Hash) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	return os.RemoveAll(fullPath)
}
