package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/nightlyone/lockfile"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
func (r *fileRepo) CreateBox(addr address.HashAddress, box int) error {
	fullPath := r.getPath(addr, getBoxAsString(box))

	return os.MkdirAll(fullPath, 0700)
}

// Returns true when the given mailbox exists in this account
func (r *fileRepo) ExistsBox(addr address.HashAddress, box int) bool {
	return r.pathExists(addr, getBoxAsString(box))
}

// Delete a given mailbox in the account
func (r *fileRepo) DeleteBox(addr address.HashAddress, box int) error {
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

	return filepath.Join(r.basePath, strAddr[:2], strAddr[2:], suffix)
}

// Retrieve a single mailbox
func (r *fileRepo) GetBoxInfo(addr address.HashAddress, box int) (*BoxInfo, error) {
	mbi := &BoxInfo{
		ID: box,
	}

	// Check number of messages in directory
	files, err := ioutil.ReadDir(r.getPath(addr, getBoxAsString(box)))
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

func (r *fileRepo) GetAllBoxes(addr address.HashAddress) ([]BoxInfo, error) {
	var list []BoxInfo

	files, err := ioutil.ReadDir(r.getPath(addr, ""))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && isBoxDir(f.Name()) {
			bi, err := r.GetBoxInfo(addr, getBoxIdFromString(f.Name()))
			if err != nil {
				continue
			}

			list = append(list, *bi)
		}
	}

	return list, nil
}

// Query messages inside mailbox
func (r *fileRepo) FetchListFromBox(addr address.HashAddress, box int, since time.Time, offset, limit int) (*MessageList, error) {
	var list = &MessageList{
		Offset:   offset,
		Limit:    limit,
		Total:    0,
		Returned: 0,
		Messages: []Message{},
	}

	files, err := ioutil.ReadDir(r.getPath(addr, getBoxAsString(box)))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		list.Total++
		if list.Returned >= list.Limit {
			continue
		}

		header, err := r.fetch(addr, filepath.Join(getBoxAsString(box), f.Name(), "header.json"))
		if err != nil {
			continue
		}
		catalog, err := r.fetch(addr, filepath.Join(getBoxAsString(box), f.Name(), "catalog"))
		if err != nil {
			continue
		}

		list.Returned++
		list.Messages = append(list.Messages, Message{
			Header: string(header),
			Catalog: catalog,
		})
	}

	return list, nil
}

// Set flag from the given message
func (r *fileRepo) SetFlag(addr address.HashAddress, box int, id string, flag string) error {
	return r.writeFlag(addr, box, id, flag, true)
}

// Unset flag from the given message
func (r *fileRepo) UnsetFlag(addr address.HashAddress, box int, id string, flag string) error {
	return r.writeFlag(addr, box, id, flag, false)
}

// Get flags from the given message
func (r *fileRepo) GetFlags(addr address.HashAddress, box int, id string) ([]string, error) {
	flags := &message.Flags{}
	err := r.fetchJSON(addr, filepath.Join(getBoxAsString(box), id, flagFile), flags)
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

func (r *fileRepo) writeFlag(addr address.HashAddress, box int, id string, flag string, addFlag bool) error {
	// Lock our flags for writing
	lockfilePath := r.getPath(addr, filepath.Join(getBoxAsString(box), id, flagFile+".lock"))
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

	return r.store(addr, filepath.Join(getBoxAsString(box), id, flagFile), data)
}

func (r *fileRepo) MoveToBox(addr address.HashAddress, srcBox, dstBox int, msgID string) error {
	srcPath := r.getPath(addr, filepath.Join(getBoxAsString(srcBox), msgID))
	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(dstBox), msgID))

	return os.Rename(srcPath, dstPath)
}

// Send a message to specific box
func (r *fileRepo) SendToBox(addr address.HashAddress, box int, msgID string) error {
	srcPath, err := message.GetPath(message.SectionProcessing, msgID, "")
	if err != nil {
		return err
	}

	dstPath := r.getPath(addr, filepath.Join(getBoxAsString(box), msgID))
	// // If we have the inbox, the message is prefixed with the current timestamp (UTC). This allows us
	// // sort on time locally and we can just fetch from a specific time (ie: fetch all messages since 20 minutes ago)
	// if box == "inbox" {
	// 	dstPath = r.getPath(addr, filepath.Join(box, fmt.Sprintf("%d-%s", time.Now().Unix(), msgID)))
	// }
	return os.Rename(srcPath, dstPath)
}

func getBoxIdFromString(dir string) int {
	if !isBoxDir(dir) {
		return 0
	}

	dir = strings.TrimPrefix(dir, "box-")
	box, err := strconv.Atoi(dir)
	if err != nil {
		return 0
	}

	return box

}

func isBoxDir(dir string) bool {
	return strings.HasPrefix(dir, "box-")
}

func getBoxAsString(box int) string {
	return fmt.Sprintf("box-%d", box)
}
