package internal

import (
	"github.com/juju/fslock"
	"io/ioutil"
	"os"
)

// WriteFileWithLock writes data by safely writing to a temp file first
func WriteFileWithLock(path string, data []byte, perm os.FileMode) error {
	// Lock the file first. Make sure we are the only one working on it
	lock := fslock.New(path + ".lock")
	err := lock.TryLock()
	if err != nil {
		return err
	}

	defer func() {
		_ = lock.Unlock()
		_ = os.Remove(path + ".lock")
	}()

	err = ioutil.WriteFile(path+".tmp", data, perm)
	if err != nil {
		return err
	}

	err = os.Rename(path+".tmp", path)
	return err
}

// ReadFileWithLock reads a file which might be locked
func ReadFileWithLock(p string) ([]byte, error) {
	// Lock vault for reading
	lock := fslock.New(p + ".lock")
	err := lock.TryLock()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(p)

	_ = lock.Unlock()
	_ = os.Remove(p + ".lock")
	if err != nil {
		return nil, err
	}

	return data, nil
}
