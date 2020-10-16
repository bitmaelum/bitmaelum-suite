// Copyright (c) 2020 BitMaelum Authors
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
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// Create a new mailbox in this account
func (r *fileRepo) CreateBox(addr hash.Hash, box int) error {
	fullPath := r.getPath(addr, getBoxAsString(box))

	return os.MkdirAll(fullPath, 0700)
}

// Returns true when the given mailbox exists in this account
func (r *fileRepo) ExistsBox(addr hash.Hash, box int) bool {
	return r.pathExists(addr, getBoxAsString(box))
}

// Delete a given mailbox in the account
func (r *fileRepo) DeleteBox(addr hash.Hash, box int) error {
	if box <= MaxMandatoryBoxID {
		return errors.New("cannot delete mandatory box")
	}

	fullPath := r.getPath(addr, getBoxAsString(box))

	err := os.RemoveAll(fullPath)
	if err != nil {
		return errors.New("cannot remove box: " + err.Error())
	}

	return err
}

// Retrieve a single mailbox
func (r *fileRepo) GetBoxInfo(addr hash.Hash, box int) (*BoxInfo, error) {
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

func (r *fileRepo) GetAllBoxes(addr hash.Hash) ([]BoxInfo, error) {
	var list []BoxInfo

	files, err := ioutil.ReadDir(r.getPath(addr, ""))
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
