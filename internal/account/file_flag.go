package account

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/nightlyone/lockfile"
	"path/filepath"
)

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
