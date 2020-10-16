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

package message

import (
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetRetryInfoFromQueue(t *testing.T) {
	config.Server.Paths.Retry = "/datastore/retry/"

	fs = afero.NewMemMapFs()
	data := `{"retry_at":"2010-01-01T12:35:56Z","last_retried_at":"2010-01-01T12:34:56Z","retries":0,"message_id":"2f9011bf-912d-4c64-91f1-bd6a99c03375"}`
	_ = afero.WriteFile(fs, "/datastore/retry/2f9011bf-912d-4c64-91f1-bd6a99c03375/.retry.json", []byte(data), 0644)
	data = `{"retry_at":"2010-01-01T12:35:56Z","last_retried_at":"2010-01-01T12:34:56Z","retries":0,"message_id":"2f9011bf-912d-4c64-91f1-bd6a99c03375"}`
	_ = afero.WriteFile(fs, "/datastore/retry/33333333-912d-4c64-91f1-bd6a99c03375/.retry.json", []byte(data), 0644)
	data = `{"retry_at":"2010-01-01T12:35:56Z","last_retried_at":"2010-01-01T12:34:56Z","retries":0,"message_id":"2f9011bf-912d-4c64-91f1-bd6a99c03375"}`
	_ = afero.WriteFile(fs, "/datastore/retry/2f9011bf-912d-4c64-91f1-bd6a99c03375/.notaretry.json", []byte(data), 0644)

	ris, err := GetRetryInfoFromQueue()
	assert.NoError(t, err)
	assert.Len(t, ris, 2)
}

func TestNewRetryInfo(t *testing.T) {
	// Assume this is the current time during tests
	timeNow = func() time.Time {
		return time.Date(2010, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	ri := NewRetryInfo("2f9011bf-912d-4c64-91f1-bd6a99c03375")

	assert.Equal(t, timeNow().Unix(), ri.LastRetriedAt.Unix())
	assert.Equal(t, "2f9011bf-912d-4c64-91f1-bd6a99c03375", ri.MsgID)
	assert.Equal(t, 0, ri.Retries)
	assert.Equal(t, timeNow().Add(60*time.Second).Unix(), ri.RetryAt.Unix())
}

func TestStoreRetryInfo(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = fs.MkdirAll("/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375", 0755)

	ri := NewRetryInfo("2f9011bf-912d-4c64-91f1-bd6a99c03375")
	err := StoreRetryInfo(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375", *ri)
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/.retry.json")
	assert.True(t, ok)

	ri2, err := GetRetryInfo(SectionIncoming, "0000000")
	assert.Error(t, err)
	assert.Nil(t, ri2)

	ri2, err = GetRetryInfo(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)
	assert.Equal(t, 0, ri2.Retries)
	assert.Equal(t, "2f9011bf-912d-4c64-91f1-bd6a99c03375", ri2.MsgID)
}
