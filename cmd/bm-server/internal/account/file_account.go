// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package account

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Create a new account for this address
func (r *fileRepo) Create(addr hash.Hash, pubKey bmcrypto.PubKey) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	err := r.fs.MkdirAll(fullPath, 0700)
	if err != nil {
		return err
	}

	// parallelize actions
	g := new(errgroup.Group)
	g.Go(func() error {
		return r.StoreKey(addr, pubKey)
	})

	g.Go(func() error {
		messagePath := r.getPath(addr, "messages")
		return r.fs.MkdirAll(messagePath, 0700)
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

	return r.fs.RemoveAll(fullPath)
}
