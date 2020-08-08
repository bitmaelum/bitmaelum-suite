package account

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	pubKeyFile = ".keys.json"
	infoFile   = ".info.json"
	flagFile   = ".flags.json"
)

// PubKeys holds a list of public keys
type PubKeys struct {
	PubKeys []string `json:"keys"`
}

type fileRepo struct {
	basePath string
}

// NewFileRepository returns a new file repository
func NewFileRepository(basePath string) Repository {
	return &fileRepo{
		basePath: basePath,
	}
}

// Store data on the given account path
func (r *fileRepo) store(addr address.HashAddress, path string, data []byte) error {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("storing file on %s", fullPath)

	return ioutil.WriteFile(fullPath, data, 0600)
}

// Check if path in account exists
func (r *fileRepo) pathExists(addr address.HashAddress, path string) bool {
	fullPath := r.getPath(addr, path)
	_, err := os.Stat(fullPath)

	return !os.IsNotExist(err)
}

// Delete path in account
func (r *fileRepo) delete(addr address.HashAddress, path string) error {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("deleting file %s", fullPath)

	return os.Remove(fullPath)
}

// Retrieve data on path in account
func (r *fileRepo) fetch(addr address.HashAddress, path string) ([]byte, error) {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("fetching file %s", fullPath)

	return ioutil.ReadFile(fullPath)
}

// Retrieves a data structure based on JSON
func (r *fileRepo) fetchJSON(addr address.HashAddress, path string, v interface{}) error {
	fullPath := r.getPath(addr, path)
	logrus.Debugf("fetching file %s", fullPath)

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

// Generate the path in account
func (r *fileRepo) getPath(addr address.HashAddress, suffix string) string {
	strAddr := strings.ToLower(addr.String())
	suffix = strings.ToLower(suffix)

	return filepath.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}
