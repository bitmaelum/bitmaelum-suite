package account

import (
    "encoding/json"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/message"
    "github.com/nightlyone/lockfile"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
)


const (
    PUBKEY_FILE = ".pubkeys.json"
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
func (r *fileRepo) StorePubKey(addr core.HashAddress, key string) error {
    // Lock our keyfile for writing
    lock, err := lockfile.New(PUBKEY_FILE + ".lock")

    err = lock.TryLock()
    if err != nil {
        return err
    }

    defer func () {
        _ = lock.Unlock()
    }()

    // Read keys
    pk := &message.Pubkeys{}
    err = r.fetchJson(addr, PUBKEY_FILE, pk)
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
    return r.store(addr, PUBKEY_FILE, data)
}

// Retrieve the public key for this account
func (r *fileRepo) FetchPubKeys(addr core.HashAddress) ([]string, error) {
    pk := &message.Pubkeys{}
    err := r.fetchJson(addr, PUBKEY_FILE, pk)
    if err != nil {
        return nil, err
    }

    return pk.PubKeys, nil
}

// Create a new mailbox in this account
func (r *fileRepo) CreateBox(addr core.HashAddress, box, description string, quota int) error {
    fullPath := r.getPath(addr, box)

    _ = os.MkdirAll(fullPath, 0700)

    mbi := message.MailBoxInfo{
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

// Retrieves a data structure based on JSON
func (r *fileRepo) fetchJson(addr core.HashAddress, path string, v interface{}) error {
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
func (r *fileRepo) getPath(addr core.HashAddress, suffix string) string {
    strAddr := strings.ToLower(addr.String())
    suffix = strings.ToLower(suffix)

    return path.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}

// Retrieve a single mailbox
func (r *fileRepo) GetBox(addr core.HashAddress, box string) (*message.MailBoxInfo, error) {
    mbi := &message.MailBoxInfo{}
    err := r.fetchJson(addr, path.Join(box, INFO_FILE), mbi)
    if err != nil {
        return nil, err
    }

    mbi.Name = box
    return mbi, nil
}

// Search for mailboxes. Use glob-patterns for querying
func (r *fileRepo) FindBox(addr core.HashAddress, query string) ([]message.MailBoxInfo, error) {
    list := []message.MailBoxInfo{}

    files, err := ioutil.ReadDir(r.getPath(addr, ""))
    if err != nil {
        return nil, err
    }

    for _, f := range files {
        matched, err := filepath.Match(query, f.Name())
        if ! matched || err != nil {
            continue
        }

        boxInfo, err := r.GetBox(addr, f.Name())
        if err != nil {
            continue
        }

        list = append(list, *boxInfo)
    }

    return list, nil
}

func (r *fileRepo) FindMessages(addr core.HashAddress, box string, offset, limit int) ([]message.MessageInfo, error) {
    list := []message.MessageInfo{}

    files, err := ioutil.ReadDir(r.getPath(addr, box))
    if err != nil {
        return nil, err
    }

    for _, f := range files {
        if ! f.IsDir() {
            continue
        }

        mi, err := r.GetMessageInfo(addr, box, f.Name())
        if err != nil {
            continue
        }
        list = append(list, *mi)
    }

    return list, nil
}

// Fetch specific mail
func (r *fileRepo) GetMessageInfo(addr core.HashAddress, box string, msgUuid string) (*message.MessageInfo, error) {

    c := &message.Catalog{}
    err := r.fetchJson(addr, path.Join(box, msgUuid, "catalog.json"), c)
    if err != nil {
        return nil, err
    }

    f := &message.Flags{}
    _ = r.fetchJson(addr, path.Join(box, msgUuid, ".flags"), f)
    
    return &message.MessageInfo{
        Flags:   *f,
        Catalog: *c,
    }, nil
}
