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

package message

import (
	"path/filepath"
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"
	config.Server.Paths.Accounts = "/datastore/accounts/"
	config.Server.Paths.Processing = "/datastore/processing/"
	config.Server.Paths.Retry = "/datastore/retry/"

	var (
		s   string
		err error
	)

	s, err = GetPath(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/datastore", "incoming", "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt"), s)

	s, err = GetPath(SectionProcessing, "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/datastore", "processing", "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt"), s)

	s, err = GetPath(SectionRetry, "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/datastore", "retry", "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt"), s)

	_, err = GetPath(5232, "2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt")
	assert.Error(t, err)

	_, err = GetPath(SectionIncoming, "../../../../../../etc/passwd", "foo.txt")
	assert.Error(t, err)
}

func TestIncomingPathExists(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/foo.txt", []byte("foobar"), 0644)

	assert.True(t, IncomingPathExists("2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt"))
	assert.False(t, IncomingPathExists("2f9011bf-912d-4c64-91f1-bd6a99c03375", "not-exist"))

	assert.False(t, IncomingPathExists("11111111-1111-1111-1111-111111111111", "foo.txt"))
	assert.False(t, IncomingPathExists("2f9011bf-912d-4c64-91f1-bd6a99c033751", ""))
}

func TestProcessQueuePathExists(t *testing.T) {
	config.Server.Paths.Processing = "/datastore/processing/"

	fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/datastore/processing/2f9011bf-912d-4c64-91f1-bd6a99c03375/foo.txt", []byte("foobar"), 0644)

	assert.True(t, ProcessQueuePathExists("2f9011bf-912d-4c64-91f1-bd6a99c03375", "foo.txt"))
	assert.False(t, ProcessQueuePathExists("2f9011bf-912d-4c64-91f1-bd6a99c03375", "not-exist"))

	assert.False(t, ProcessQueuePathExists("11111111-1111-1111-1111-111111111111", "foo.txt"))
	assert.False(t, ProcessQueuePathExists("2f9011bf-912d-4c64-91f1-bd6a99c033751", ""))
}
