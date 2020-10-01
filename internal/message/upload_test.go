package message

import (
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetMessageHeader(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)
	_ = afero.WriteFile(fs, "/datastore/incoming/33333333-912d-4c64-91f1-bd6a99c03375/header.json", []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"), 0644)

	h, err := GetMessageHeader(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)
	assert.Equal(t, "000000000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc", h.From.Addr.String())
	assert.Equal(t, uint64(983), h.Catalog.Size)

	// not existing
	h, err = GetMessageHeader(SectionIncoming, "11111111-1111-1111-1111-111111111111")
	assert.Error(t, err)
	assert.Nil(t, h)

	// invalid
	h, err = GetMessageHeader(SectionIncoming, "33333333-912d-4c64-91f1-bd6a99c03375")
	assert.Error(t, err)
	assert.Nil(t, h)

}

func TestRemoveMessage(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/test.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)
	_ = afero.WriteFile(fs, "/datastore/incoming/33333333-912d-4c64-91f1-bd6a99c03375/header.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)

	err := RemoveMessage(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)

	var ok bool
	ok, _ = afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json")
	assert.False(t, ok)
	ok, _ = afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/test.json")
	assert.False(t, ok)
	ok, _ = afero.DirExists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.False(t, ok)

	ok, _ = afero.Exists(fs, "/datastore/incoming/33333333-912d-4c64-91f1-bd6a99c03375/header.json")
	assert.True(t, ok)
}

func TestMoveMessage(t *testing.T) {
	t.Skip("We must skip move message as afero memmap system does not work correctly with message dir changes")

	config.Server.Paths.Incoming = "/datastore/incoming/"
	config.Server.Paths.Retry = "/datastore/retry/"

	fs = afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)
	_ = afero.WriteFile(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/test.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)
	_ = afero.WriteFile(fs, "/datastore/incoming/33333333-912d-4c64-91f1-bd6a99c03375/header.json", internal.ReadTestFile("../../testdata/header-001.json"), 0644)

	err := MoveMessage(SectionRetry, SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.Error(t, err)

	err = MoveMessage(SectionIncoming, SectionRetry, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)

	var ok bool
	ok, err = afero.Exists(fs, "/datastore/retry/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json")
	assert.NoError(t, err)
	assert.True(t, ok)
	ok, _ = afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestStoreHeader(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = fs.MkdirAll("/datastore/incoming", 0755)

	h := &Header{}
	_ = internal.ReadHeader("../../testdata/header-001.json", h)
	err := StoreHeader("2f9011bf-912d-4c64-91f1-bd6a99c03375", h)
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/header.json")
	assert.NoError(t, err)
	assert.True(t, ok)

	h2, err := GetMessageHeader(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)
	assert.Equal(t, "000000000000000000018f66a0f3591a883f2b9cc3e95a497e7cf9da1071b4cc", h2.From.Addr.String())
}

func TestStoreCatalog(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = fs.MkdirAll("/datastore/incoming", 0755)

	err := StoreCatalog("2f9011bf-912d-4c64-91f1-bd6a99c03375", strings.NewReader("foobar"))
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/catalog")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, _ = afero.FileContainsBytes(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/catalog", []byte("foobar"))
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestStoreBlock(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = fs.MkdirAll("/datastore/incoming", 0755)

	err := StoreBlock("2f9011bf-912d-4c64-91f1-bd6a99c03375", "7a355267-b2eb-4377-a210-63d8835faffe", strings.NewReader("foobar"))
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/7a355267-b2eb-4377-a210-63d8835faffe")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, _ = afero.FileContainsBytes(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/7a355267-b2eb-4377-a210-63d8835faffe", []byte("foobar"))
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestStoreAttachment(t *testing.T) {
	config.Server.Paths.Incoming = "/datastore/incoming/"

	fs = afero.NewMemMapFs()
	_ = fs.MkdirAll("/datastore/incoming", 0755)

	err := StoreAttachment("2f9011bf-912d-4c64-91f1-bd6a99c03375", "7a355267-b2eb-4377-a210-63d8835faffe", strings.NewReader("foobar"))
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/7a355267-b2eb-4377-a210-63d8835faffe")
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, _ = afero.FileContainsBytes(fs, "/datastore/incoming/2f9011bf-912d-4c64-91f1-bd6a99c03375/7a355267-b2eb-4377-a210-63d8835faffe", []byte("foobar"))
	assert.NoError(t, err)
	assert.True(t, ok)

	f, err := GetFiles(SectionIncoming, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.NoError(t, err)
	assert.Len(t, f, 1)
	assert.Equal(t, "7a355267-b2eb-4377-a210-63d8835faffe", f[0].ID)
	assert.Contains(t, f[0].Path, "2f9011bf-912d-4c64-91f1-bd6a99c03375")
	assert.Contains(t, f[0].Path, "7a355267-b2eb-4377-a210-63d8835faffe")
}
