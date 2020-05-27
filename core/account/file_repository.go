package account

import (
    "encoding/json"
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

func NewFileRepository(basePath string) Repository {
    return &fileRepo{
        basePath: basePath,
    }
}

func (r *fileRepo) Create(hash string) error {
    fullPath := r.getPath(hash, "")
    logrus.Debugf("creating hash directory %s", fullPath)

    return os.MkdirAll(fullPath, 0700)
}

func (r *fileRepo) Exists(hash string) bool {
    return r.pathExists(hash, "")
}

func (r *fileRepo) StorePubKey(hash string, data []byte) error {
    return r.store(hash, PUBKEY_FILE, data)
}

func (r *fileRepo) FetchPubKey(hash string) ([]byte, error) {
    return r.fetch(hash, PUBKEY_FILE)
}

func (r *fileRepo) CreateBox(hash string, box string, description string, quota int) error {
    fullPath := r.getPath(hash, box)

    _ = os.MkdirAll(fullPath, 0700)

    mbi := MailBoxInfo{
        Description: description,
        Quota: quota,
    }

    data, _ := json.MarshalIndent(mbi, "", " ")
    return r.store(hash, path.Join(box, INFO_FILE), data)
}

func (r *fileRepo) ExistsBox(hash string, box string) bool {
    return r.pathExists(hash, box)
}

func (r *fileRepo) DeleteBox(hash string, box string) error {
    // @TODO: not yet implemented
    panic("not implemented yet")
}


func (r *fileRepo) store(hash string, path string, data []byte) error {
    fullPath := r.getPath(hash, path)
    logrus.Debugf("storing file on %s", fullPath)

    return ioutil.WriteFile(fullPath, data, 0600)
}

func (r *fileRepo) pathExists(hash string, path string) bool {
    fullPath := r.getPath(hash, path)
    _, err := os.Stat(fullPath)

    return ! os.IsNotExist(err)
}

func (r *fileRepo) delete(hash string, path string) error {
    fullPath := r.getPath(hash, path)
    logrus.Debugf("deleting file %s", fullPath)

    return os.Remove(fullPath)
}

func (r *fileRepo) fetch(hash string, path string) ([]byte, error) {
    fullPath := r.getPath(hash, path)
    logrus.Debugf("fetching file %s", fullPath)

    return ioutil.ReadFile(fullPath)
}

func (r *fileRepo) getPath(hash string, suffix string) string {
    hash = strings.ToLower(hash)
    suffix = strings.ToLower(suffix)

    return path.Join(r.basePath, hash[:2], hash[2:], suffix)
}
