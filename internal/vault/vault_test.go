package vault

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	bmtest "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	// inmemory vault
	v, err := New("", []byte{})
	assert.NoError(t, err)
	assert.NotNil(t, v)

	fs = afero.NewMemMapFs()

	// open/create new vault
	v, err = New("/my/dir/vault.json", []byte("secret"))
	assert.NoError(t, err)
	assert.NotNil(t, v)
	ok, _ := afero.Exists(fs, "/my/dir/vault.json")
	assert.True(t, ok)

	// Reopen again
	v, err = New("/my/dir/vault.json", []byte("secret"))
	assert.NoError(t, err)
	assert.NotNil(t, v)

	// refuse to overwrite vault with empty vault
	err = v.WriteToDisk()
	assert.Errorf(t, err, "vault seems to have invalid data. Refusing to overwrite the current vault")

	privKey, pubKey, _ := bmtest.ReadTestKey("../../testdata/key-1.json")

	acc := &internal.AccountInfo{
		Default:   false,
		Address:   "foobar!",
		Name:      "Foo Bar",
		Settings:  nil,
		PrivKey:   *privKey,
		PubKey:    *pubKey,
		Pow:       &proofofwork.ProofOfWork{},
		RoutingID: "12345678",
	}
	v.AddAccount(*acc)

	// Write to disk works when we have at least one account
	err = v.WriteToDisk()
	assert.NoError(t, err)

	// Check if the backup exists
	ok, _ = afero.Exists(fs, "/my/dir/vault.json.backup")
	assert.True(t, ok)

	// Open vault with wrong password
	v, err = New("/my/dir/vault.json", []byte("incorrect password"))
	assert.Errorf(t, err, "incorrect password")
	assert.Nil(t, v)

	// Open vault with correct password
	v, err = New("/my/dir/vault.json", []byte("secret"))
	assert.NoError(t, err)
	assert.Len(t, v.Store.Accounts, 1)
}

func Test_FindShortRoutingId(t *testing.T) {
	var acc internal.AccountInfo

	v, _ := New("", []byte{})

	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780000"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780001"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780002"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "154353535335"}
	v.AddAccount(acc)

	assert.Equal(t, "154353535335", v.FindShortRoutingID("154"))
	assert.Equal(t, "154353535335", v.FindShortRoutingID("15435"))
	assert.Equal(t, "", v.FindShortRoutingID("12345"))
	assert.Equal(t, "", v.FindShortRoutingID("1"))
}
