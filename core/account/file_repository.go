package account

import (
    "encoding/json"
    "github.com/jaytaph/mailv2/core"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
    "strings"
)

// Structure of the .info.json file
type MailBoxInfo struct {
    Description string `json:"description"`
    Quota int `json:"quota"`
}

const (
    PUBKEY_FILE = ".pubkey"
    INFO_FILE = ".info.json"
)

type fileRepo struct {
    basePath string
}

// Return a new file repository
func NewFileRepository(basePath string) Repository {
    return &fileRepo{
        basePath: basePath,
    }
}

// Create a new account for this address
func (r *fileRepo) Create(addr core.HashAddress) error {
    fullPath := r.getPath(addr, "")
    logrus.Debugf("creating hash directory %s", fullPath)

    return os.MkdirAll(fullPath, 0700)
}

// Returns true when the given account for this address exists
func (r *fileRepo) Exists(addr core.HashAddress) bool {
    return r.pathExists(addr, "")
}

// Store the public key for this account
func (r *fileRepo) StorePubKey(addr core.HashAddress, data []byte) error {
    return r.store(addr, PUBKEY_FILE, data)
}

// Retrieve the public key for this account
func (r *fileRepo) FetchPubKey(addr core.HashAddress) ([]byte, error) {
    return r.fetch(addr, PUBKEY_FILE)
}

// Create a new mailbox in this account
func (r *fileRepo) CreateBox(addr core.HashAddress, box, description string, quota int) error {
    fullPath := r.getPath(addr, box)

    _ = os.MkdirAll(fullPath, 0700)

    mbi := MailBoxInfo{
        Description: description,
        Quota: quota,
    }

    data, _ := json.MarshalIndent(mbi, "", " ")
    return r.store(addr, path.Join(box, INFO_FILE), data)
}

// Returns true when the given mailbox exists in this account
func (r *fileRepo) ExistsBox(addr core.HashAddress, box string) bool {
    return r.pathExists(addr, box)
}

// Delete a given mailbox in the account
func (r *fileRepo) DeleteBox(addr core.HashAddress, box string) error {
    // @TODO: not yet implemented
    panic("not implemented yet")
}

// Store data on the given account path
func (r *fileRepo) store(addr core.HashAddress, path string, data []byte) error {
    fullPath := r.getPath(addr, path)
    logrus.Debugf("storing file on %s", fullPath)

    return ioutil.WriteFile(fullPath, data, 0600)
}

// Check if path in account exists
func (r *fileRepo) pathExists(addr core.HashAddress, path string) bool {
    fullPath := r.getPath(addr, path)
    _, err := os.Stat(fullPath)

    return ! os.IsNotExist(err)
}

// Delete path in account
func (r *fileRepo) delete(addr core.HashAddress, path string) error {
    fullPath := r.getPath(addr, path)
    logrus.Debugf("deleting file %s", fullPath)

    return os.Remove(fullPath)
}

// Retrieve data on path in account
func (r *fileRepo) fetch(addr core.HashAddress, path string) ([]byte, error) {
    fullPath := r.getPath(addr, path)
    logrus.Debugf("fetching file %s", fullPath)

    return ioutil.ReadFile(fullPath)
}

// Generate the path in account
func (r *fileRepo) getPath(addr core.HashAddress, suffix string) string {
    strAddr := strings.ToLower(addr.String())
    suffix = strings.ToLower(suffix)

    return path.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}
