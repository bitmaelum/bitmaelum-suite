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

package vault

import (
	"testing"

	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	fs = afero.NewMemMapFs()

	// inmemory vault
	v := New()
	assert.NotNil(t, v)

	// open/create new vault
	v, err := Create("/my/dir/vault.json", "secret")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	ok, _ := afero.Exists(fs, "/my/dir/vault.json")
	assert.True(t, ok)

	// Reopen again
	v, err = Open("/my/dir/vault.json", "secret")
	assert.NoError(t, err)
	assert.NotNil(t, v)

	privKey, pubKey, _ := testing2.ReadTestKey("../../testdata/key-1.json")

	addr, _ := address.NewAddress("foobar!")
	acc := &AccountInfo{
		Address:   addr,
		Name:      "Foo Bar",
		Settings:  nil,
		PrivKey:   *privKey,
		PubKey:    *pubKey,
		Pow:       &proofofwork.ProofOfWork{},
		RoutingID: "12345678",
	}
	v.AddAccount(*acc)

	// Write to disk
	err = v.Persist()
	assert.NoError(t, err)

	// Check if the backup exists
	ok, _ = afero.Exists(fs, "/my/dir/vault.json.backup")
	assert.True(t, ok)

	// Open vault with wrong password
	v, err = Open("/my/dir/vault.json", "incorrect password")
	assert.Errorf(t, err, "incorrect password")
	assert.Nil(t, v)

	// Open vault with correct password
	v, err = Open("/my/dir/vault.json", "secret")
	assert.NoError(t, err)
	assert.Len(t, v.Store.Accounts, 1)
}

func TestFindShortRoutingId(t *testing.T) {
	var acc AccountInfo

	v := New()

	addr, _ := address.NewAddress("example!")
	acc = AccountInfo{Address: addr, RoutingID: "123456780000"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: addr, RoutingID: "123456780001"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: addr, RoutingID: "123456780002"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: addr, RoutingID: "154353535335"}
	v.AddAccount(acc)

	assert.Equal(t, "154353535335", v.FindShortRoutingID("154"))
	assert.Equal(t, "154353535335", v.FindShortRoutingID("15435"))
	assert.Equal(t, "", v.FindShortRoutingID("12345"))
	assert.Equal(t, "", v.FindShortRoutingID("1"))
}

func TestNewPersistent(t *testing.T) {
	v := NewPersistent("/v1.json", "foobar")
	assert.NotNil(t, v)

	assert.Equal(t, "foobar", v.password)
	assert.Equal(t, "/v1.json", v.path)
}

func TestVaultChangePassword(t *testing.T) {
	v := New()
	assert.NotNil(t, v)

	v.SetPassword("foobar")
	assert.Equal(t, "foobar", v.password)

	v.SetPassword("secret")
	assert.Equal(t, "secret", v.password)
}

func TestExisting(t *testing.T) {
	fs = afero.NewMemMapFs()

	// Cant open vaults
	v, err := Open("/v1.json", "foo")
	assert.Error(t, err)
	assert.Nil(t, v)
	v, err = Open("/v2.json", "bar")
	assert.Error(t, err)
	assert.Nil(t, v)

	// Create vaults
	_, err = Create("/v1.json", "foo")
	assert.NoError(t, err)
	_, err = Create("/v2.json", "bar")
	assert.NoError(t, err)

	// Open vaults
	v, err = Open("/v1.json", "foo")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	v, err = Open("/v2.json", "bar")
	assert.NoError(t, err)
	assert.NotNil(t, v)

	// Change password of vault
	v.SetPassword("anotherpass")
	err = v.Persist()
	assert.NoError(t, err)

	// Try open vault again
	v, err = Open("/v2.json", "bar")
	assert.Error(t, err)
	assert.Nil(t, v)
	v, err = Open("/v2.json", "anotherpass")
	assert.NoError(t, err)
	assert.NotNil(t, v)

}
