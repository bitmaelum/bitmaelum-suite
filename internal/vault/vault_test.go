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

	privKey, pubKey, _ := testing2.ReadTestKey("../../testdata/key-1.json")

	addr, _ := address.NewAddress("foobar!")
	acc := &AccountInfo{
		Default:   false,
		Address:   *addr,
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

func TestFindShortRoutingId(t *testing.T) {
	var acc AccountInfo

	v, _ := New("", []byte{})

	addr, _ := address.NewAddress("example!")
	acc = AccountInfo{Address: *addr, RoutingID: "123456780000"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: *addr, RoutingID: "123456780001"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: *addr, RoutingID: "123456780002"}
	v.AddAccount(acc)
	acc = AccountInfo{Address: *addr, RoutingID: "154353535335"}
	v.AddAccount(acc)

	assert.Equal(t, "154353535335", v.FindShortRoutingID("154"))
	assert.Equal(t, "154353535335", v.FindShortRoutingID("15435"))
	assert.Equal(t, "", v.FindShortRoutingID("12345"))
	assert.Equal(t, "", v.FindShortRoutingID("1"))
}

func TestVaultChangePassword(t *testing.T) {
	v, _ := New("", []byte("foobar"))
	assert.Equal(t, []byte("foobar"), v.password)

	v.ChangePassword("secret")
	assert.Equal(t, []byte("secret"), v.password)
}

func TestGetAccountOrDefault(t *testing.T) {
	var acc *AccountInfo

	v, _ := New("", []byte{})
	addr, _ := address.NewAddress("example1!")
	acc = &AccountInfo{Address: *addr, RoutingID: "123456780000"}
	v.AddAccount(*acc)
	addr, _ = address.NewAddress("example2!")
	acc = &AccountInfo{Address: *addr, RoutingID: "123456780001"}
	v.AddAccount(*acc)
	addr, _ = address.NewAddress("example3!")
	acc = &AccountInfo{Address: *addr, RoutingID: "123456780002", Default: true}
	v.AddAccount(*acc)
	addr, _ = address.NewAddress("example4!")
	acc = &AccountInfo{Address: *addr, RoutingID: "154353535335"}
	v.AddAccount(*acc)

	acc = GetAccountOrDefault(v, "")
	assert.Equal(t, "example3!", acc.Address)

	acc = GetAccountOrDefault(v, "example2!")
	assert.Equal(t, "example2!", acc.Address)
}
