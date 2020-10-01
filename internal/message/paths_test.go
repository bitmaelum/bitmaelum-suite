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
