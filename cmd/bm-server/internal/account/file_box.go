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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/afero"
)

var (
	errCantDeleteMandatoryBox = errors.New("cannot delete mandatory box")
	errCannotDeleteBox        = errors.New("cannot remove box")
)

// Create a new mailbox in this account
func (r *fileRepo) CreateBox(addr hash.Hash, box int) error {
	fullPath := r.getPath(addr, getBoxAsString(box))

	return r.fs.MkdirAll(fullPath, 0700)
}

// Returns true when the given mailbox exists in this account
func (r *fileRepo) ExistsBox(addr hash.Hash, box int) bool {
	return r.pathExists(addr, getBoxAsString(box))
}

// Delete a given mailbox in the account
func (r *fileRepo) DeleteBox(addr hash.Hash, box int) error {
	if box <= MaxMandatoryBoxID {
		return errCantDeleteMandatoryBox
	}

	fullPath := r.getPath(addr, getBoxAsString(box))

	err := r.fs.RemoveAll(fullPath)
	if err != nil {
		return errCannotDeleteBox
	}

	return err
}

// Retrieve a single mailbox
func (r *fileRepo) GetBoxInfo(addr hash.Hash, box int) (*BoxInfo, error) {
	mbi := &BoxInfo{
		ID: box,
	}

	// Check number of messages in directory
	files, err := afero.ReadDir(r.fs, r.getPath(addr, getBoxAsString(box)))
	if err != nil {
		mbi.Total = 0
	} else {
		for _, file := range files {
			// File should be symlink. In the old days, it was a directory instead
			if file.IsDir() || file.Mode()&os.ModeSymlink != 0 {
				mbi.Total++
				mbi.Messages = append(mbi.Messages, file.Name())
			}
		}
	}

	return mbi, nil
}

func (r *fileRepo) GetAllBoxes(addr hash.Hash) ([]BoxInfo, error) {
	var list []BoxInfo

	files, err := afero.ReadDir(r.fs, r.getPath(addr, ""))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() && isBoxDir(f.Name()) {
			bi, err := r.GetBoxInfo(addr, getBoxIDFromString(f.Name()))
			if err != nil {
				continue
			}

			list = append(list, *bi)
		}
	}

	return list, nil
}

func getBoxIDFromString(dir string) int {
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
