// Copyright (c) 2021 BitMaelum Authors
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
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	keysFile         = ".keys.json"
	organisationFile = ".organisation.json"
)

// PubKeys holds a list of public keys
type PubKeys struct {
	PubKeys []bmcrypto.PubKey `json:"keys"`
}

type fileRepo struct {
	basePath string
	fs       afero.Fs
	canLock  bool // Can the filesystem do file locking?
}

// NewFileRepository returns a new file repository
func NewFileRepository(basePath string) Repository {
	return &fileRepo{
		basePath: basePath,
		fs:       afero.NewOsFs(),
		canLock:  true,
	}
}

// SetFileSystem sets the filesystem. This is mostly used for setting an afero mock file system so we can easily mock
// this system
func (r *fileRepo) SetFileSystem(fs afero.Fs, canLock bool) {
	r.fs = fs
	r.canLock = canLock
}

// lockFile will lock the given file and returns it IF the filesystem currently initialized allows locking. It will return
// a nil lockfile (which is valid) if no locking is possible.
func (r *fileRepo) lockFile(p string) (*lockfile.Lockfile, error) {
	if !r.canLock {
		return nil, nil
	}

	l, err := lockfile.New(p)
	return &l, err
}

// unLockFile unlocks the filelock if one is given
func (r *fileRepo) unLockFile(l *lockfile.Lockfile) error {
	if l != nil {
		return l.Unlock()
	}

	return nil
}

// tryLockFile tries the filelock if one is given
func (r *fileRepo) tryLockFile(l *lockfile.Lockfile) error {
	if l != nil {
		return l.TryLock()
	}

	return nil
}

// Store data on the given account path
func (r *fileRepo) store(addr hash.Hash, path string, data []byte) error {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("storing file on %s", fullPath)

	return afero.WriteFile(r.fs, fullPath, data, 0600)
}

// Check if path in account exists
func (r *fileRepo) pathExists(addr hash.Hash, path string) bool {
	fullPath := r.getPath(addr, path)
	_, err := r.fs.Stat(fullPath)

	return !os.IsNotExist(err)
}

// // Delete path in account
// func (r *fileRepo) delete(addr hash.Hash, path string) error {
// 	fullPath := r.getPath(addr, path)
// 	logrus.Debugf("deleting file %s", fullPath)
//
// 	return os.Remove(fullPath)
// }

// Retrieve data on path in account
func (r *fileRepo) fetch(addr hash.Hash, path string) ([]byte, error) {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("fetching file %s", fullPath)

	b, err := afero.ReadFile(r.fs, fullPath)
	if err != nil {
		logrus.Tracef("file: cannot read file: %s", fullPath)
		return nil, err
	}

	return b, nil
}

func (r *fileRepo) fetchReader(addr hash.Hash, path string) (rdr io.ReadCloser, size int64, err error) {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("fetching file reader %s", fullPath)

	f, err := r.fs.Open(fullPath)
	if err != nil {
		return nil, 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	return f, info.Size(), nil
}

// Retrieves a data structure based on JSON
func (r *fileRepo) fetchJSON(addr hash.Hash, path string, v interface{}) error {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("fetching file %s", fullPath)

	data, err := afero.ReadFile(r.fs, fullPath)
	if err != nil {
		logrus.Tracef("file: cannot read file: %s", fullPath)
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		logrus.Tracef("file: cannot unmarshal file: %s", fullPath)
		return err
	}

	return nil
}

// Generate the path in account
func (r *fileRepo) getPath(addr hash.Hash, suffix string) string {
	strAddr := strings.ToLower(addr.String())
	suffix = strings.ToLower(suffix)

	if len(strAddr) < 2 {
		// @TODO: we should probably not panic here, but not sure how to deal with this issue
		logrus.Panicf("Path seems wrong: '%s'", strAddr)
	}

	return filepath.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}
