// Copyright (c) 2022 BitMaelum Authors
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

package store

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestBoltStorage(t *testing.T) {
	var (
		ok  bool
		err error
	)

	rand.Seed(time.Now().UnixNano())

	acc1 := hash.New("foo!")
	acc2 := hash.New("bar!")

	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	})

	// Unfortunately, boltdb cannot be used with afero
	path := filepath.Join(os.TempDir(), fmt.Sprintf("store-%d", rand.Int31()))
	_ = os.MkdirAll(filepath.Join(path, acc1.String()[:2], acc1.String()[2:]), 0755)

	defer func() {
		_ = os.RemoveAll(path)
	}()

	b := NewBoltRepository(path)
	assert.NotNil(t, b)

	err = b.OpenDb(acc1)
	assert.NoError(t, err)

	// Initially, only root should be present
	ok = b.HasEntry(acc1, makeHash(acc1, "/"))
	assert.True(t, ok)
	ok = b.HasEntry(acc1, makeHash(acc1, "/something"))
	assert.False(t, ok)

	// Incorrect hash
	ok = b.HasEntry(acc1, makeHash(acc2, "/"))
	assert.False(t, ok)

	// Get root entry
	entry, err := b.GetEntry(acc1, makeHash(acc1, "/"), false, time.Time{})
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.False(t, entry.HasChildren)
	assert.Len(t, entry.Children, 0)

	// Add entry
	entry2 := mockEntry(acc1, "foobar", "/contacts", "/")
	err = b.SetEntry(acc1, entry2)
	assert.NoError(t, err)

	ok = b.HasEntry(acc1, makeHash(acc1, "/"))
	assert.True(t, ok)
	ok = b.HasEntry(acc1, makeHash(acc1, "/contacts"))
	assert.True(t, ok)

	entry, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.NotNil(t, entry)
	assert.Equal(t, "9f198242afd0a2660077b05c90c4aad8807b381f8e1af89e556c9a0e0e66331d", entry.Path.String())
	assert.Equal(t, "94723340d93b27ca21384fa64db760e10ee2382a3ded94f1e4243bacc24825e6", entry.Parent.String())
	assert.Equal(t, []byte("foobar"), entry.Data)

	// Create entry without correct path
	entry2 = mockEntry(acc1, "foobar", "/path/not/exist/item", "/path/not/exist")
	err = b.SetEntry(acc1, entry2)
	assert.Error(t, err)
}

func mockEntry(acc hash.Hash, data string, path string, parent string) EntryType {
	entry := NewEntry([]byte(data))

	entry.Path = makeHash(acc, path)
	p := makeHash(acc, parent)
	entry.Parent = &p

	return entry
}

func TestTimePropagation(t *testing.T) {
	acc1 := hash.New("foo!")
	path, b := setup(t, acc1)

	defer func() {
		_ = os.RemoveAll(path)
	}()

	defer func() {
		_ = os.RemoveAll(path)
	}()

	entry2, err := b.GetEntry(acc1, makeHash(acc1, "/contacts"), true, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 3)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/1"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/2"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	// New entry

	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 05, 05, 12, 34, 56, 0, time.UTC)
	})

	entry := mockEntry(acc1, "latest entry", "/contacts/7", "/contacts")
	err = b.SetEntry(acc1, entry)
	assert.NoError(t, err)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), true, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 4)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/1"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/7"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1273062896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/2"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1273062896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1273062896), entry2.Timestamp)

	// Update entry

	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 8, 8, 12, 34, 56, 0, time.UTC)
	})

	entry = mockEntry(acc1, "update entry", "/contacts/2", "/contacts")
	err = b.SetEntry(acc1, entry)
	assert.NoError(t, err)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), true, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 4)
	assert.True(t, entry2.HasChildren)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 0)
	assert.True(t, entry2.HasChildren)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/1"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/7"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1273062896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/2"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1281270896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1281270896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1281270896), entry2.Timestamp)
}

func TestRemoveEntries(t *testing.T) {
	acc1 := hash.New("foo!")
	path, b := setup(t, acc1)

	defer func() {
		_ = os.RemoveAll(path)
	}()

	entry2, err := b.GetEntry(acc1, makeHash(acc1, "/contacts"), true, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 3)
	assert.True(t, entry2.HasChildren)

	// Remove entry

	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 8, 8, 12, 34, 56, 0, time.UTC)
	})

	// Cannot remove collections
	err = b.RemoveEntry(acc1, makeHash(acc1, "/contacts"), true)
	assert.Error(t, err)
	err = b.RemoveEntry(acc1, makeHash(acc1, "/"), true)
	assert.Error(t, err)

	// Remove second entry
	err = b.RemoveEntry(acc1, makeHash(acc1, "/contacts/2"), true)
	assert.NoError(t, err)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), true, time.Time{})
	assert.NoError(t, err)
	assert.Len(t, entry2.Children, 2)
	assert.True(t, entry2.HasChildren)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/1"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	ok := b.HasEntry(acc1, makeHash(acc1, "/contacts/2"))
	assert.False(t, ok)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts/3"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1262349296), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/contacts"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1281270896), entry2.Timestamp)

	entry2, err = b.GetEntry(acc1, makeHash(acc1, "/"), false, time.Time{})
	assert.NoError(t, err)
	assert.Equal(t, int64(1281270896), entry2.Timestamp)
}

func setup(t *testing.T, acc hash.Hash) (string, Repository) {
	var (
		err error
	)

	rand.Seed(time.Now().UnixNano())

	internal.SetMockTime(func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	})

	// Unfortunately, boltdb cannot be used with afero
	path := filepath.Join(os.TempDir(), fmt.Sprintf("store-%d", rand.Int31()))
	_ = os.MkdirAll(filepath.Join(path, acc.String()[:2], acc.String()[2:]), 0755)

	b := NewBoltRepository(path)
	assert.NotNil(t, b)

	err = b.OpenDb(acc)
	assert.NoError(t, err)

	entry := mockEntry(acc, "contact list", "/contacts", "/")
	err = b.SetEntry(acc, entry)
	assert.NoError(t, err)

	entry = mockEntry(acc, "john doe", "/contacts/1", "/contacts")
	err = b.SetEntry(acc, entry)
	assert.NoError(t, err)

	entry = mockEntry(acc, "foo bar", "/contacts/2", "/contacts")
	err = b.SetEntry(acc, entry)
	assert.NoError(t, err)

	entry = mockEntry(acc, "jane austin", "/contacts/3", "/contacts")
	err = b.SetEntry(acc, entry)
	assert.NoError(t, err)

	return path, b
}

func makeHash(account hash.Hash, path string) hash.Hash {
	return hash.New(account.String() + path)
}
