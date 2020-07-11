package account

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

// Create a new account for this address
func (r *fileRepo) Create(addr address.HashAddress) error {
	fullPath := r.getPath(addr, "")
	logrus.Debugf("creating hash directory %s", fullPath)

	return os.MkdirAll(fullPath, 0700)
}

// Returns true when the given account for this address exists
func (r *fileRepo) Exists(addr address.HashAddress) bool {
	return r.pathExists(addr, "")
}

// Store the public key for this account
func (r *fileRepo) StorePubKey(addr address.HashAddress, key string) error {
	// Lock our keyfile for writing
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

// Retrieve the public key for this account
func (r *fileRepo) FetchPubKeys(addr address.HashAddress) ([]string, error) {
	pk := &PubKeys{}
	err := r.fetchJSON(addr, pubKeyFile, pk)
	if err != nil {
		return nil, err
	}

	return pk.PubKeys, nil
}

// Create a new mailbox in this account
func (r *fileRepo) CreateBox(addr address.HashAddress, box, name, description string, quota int) error {
	fullPath := r.getPath(addr, box)

	_ = os.MkdirAll(fullPath, 0700)

	mbi := message.MailBoxInfo{
		Name:        name,
		Description: description,
		Quota:       quota,
	}

	data, _ := json.MarshalIndent(mbi, "", " ")
	return r.store(addr, path.Join(box, infoFile), data)
}

// Returns true when the given mailbox exists in this account
func (r *fileRepo) ExistsBox(addr address.HashAddress, box string) bool {
	return r.pathExists(addr, box)
}

// Delete a given mailbox in the account
func (r *fileRepo) DeleteBox(addr address.HashAddress, box string) error {
	// @TODO: not yet implemented
	return errors.New("not implemented yet")
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

	return path.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}

// Retrieve a single mailbox
func (r *fileRepo) GetBox(addr address.HashAddress, box string) (*message.MailBoxInfo, error) {
	mbi := &message.MailBoxInfo{}
	mbi.Name = box

	// Fetch information from .info file
	err := r.fetchJSON(addr, path.Join(box, infoFile), mbi)
	if err != nil {
		return nil, err
	}

	// Check number of messages in directory
	files, err := ioutil.ReadDir(r.getPath(addr, box))
	if err != nil {
		mbi.Total = 0
	} else {
		for _, file := range files {
			if file.IsDir() {
				mbi.Total++
			}
		}
	}

	return mbi, nil
}

// Search for mailboxes. Use glob-patterns for querying
func (r *fileRepo) FindBox(addr address.HashAddress, query string) ([]message.MailBoxInfo, error) {
	var list []message.MailBoxInfo

	files, err := ioutil.ReadDir(r.getPath(addr, ""))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		matched, err := filepath.Match(query, f.Name())
		if !matched || err != nil {
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

// Query messages inside mailbox
func (r *fileRepo) FetchListFromBox(addr address.HashAddress, box string, offset, limit int) ([]message.List, error) {
	var list []message.List

	files, err := ioutil.ReadDir(r.getPath(addr, box))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		flags := &message.Flags{}
		_ = r.fetchJSON(addr, path.Join(box, f.Name(), flagFile), flags)

		msg := message.List{
			ID:    f.Name(),
			Dt:    f.ModTime().Format(time.RFC3339),
			Flags: flags.Flags,
		}

		list = append(list, msg)
	}

	return list, nil
}

// Set flag from the given message
func (r *fileRepo) SetFlag(addr address.HashAddress, box string, id string, flag string) error {
	return r.writeFlag(addr, box, id, flag, true)
}

// Unset flag from the given message
func (r *fileRepo) UnsetFlag(addr address.HashAddress, box string, id string, flag string) error {
	return r.writeFlag(addr, box, id, flag, false)
}

// Get flags from the given message
func (r *fileRepo) GetFlags(addr address.HashAddress, box string, id string) ([]string, error) {
	flags := &message.Flags{}
	err := r.fetchJSON(addr, path.Join(box, id, flagFile), flags)
	if err != nil {
		return nil, err
	}

	return flags.Flags, err
}

// Remove element from slice
func remove(slice []string, item string) []string {
	idx, err := find(slice, item)
	if err != nil {
		return slice
	}

	return append(slice[:idx], slice[idx+1:]...)
}

// Find element in slice
func find(slice []string, item string) (int, error) {
	for i, n := range slice {
		if item == n {
			return i, nil
		}
	}

	return 0, errors.New("not found")
}

func (r *fileRepo) writeFlag(addr address.HashAddress, box string, id string, flag string, addFlag bool) error {
	// Lock our flags for writing
	lockfilePath := r.getPath(addr, path.Join(box, id, flagFile+".lock"))
	lockfilePath, err := filepath.Abs(lockfilePath)
	if err != nil {
		return err
	}
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

	// Get flags
	flags, err := r.GetFlags(addr, box, id)
	if err != nil {
		return err
	}

	// We remove the flag first. This also takes care of duplicate flags
	flags = remove(flags, flag)

	// Add flag if needed.
	if addFlag {
		flags = append(flags, flag)
	}

	// Save flags back
	ft := &message.Flags{
		Flags: flags,
	}
	data, err := json.MarshalIndent(ft, "", "  ")
	if err != nil {
		return err
	}

	return r.store(addr, path.Join(box, id, flagFile), data)
}

func (r *fileRepo) MoveToBox(addr address.HashAddress, srcBox, dstBox, msgID string) error {
	srcPath := r.getPath(addr, path.Join(srcBox, msgID))
	dstPath := r.getPath(addr, path.Join(dstBox, msgID))

	return os.Rename(srcPath, dstPath)
}

// Send a message to specific inbox
func (r *fileRepo) SendToBox(addr address.HashAddress, box, msgID string) error {
	srcPath, err := message.GetPath(message.SectionProcessing, msgID, "")
	if err != nil {
		return err
	}

	dstPath := r.getPath(addr, path.Join(box, msgID))

	return os.Rename(srcPath, dstPath)
}
