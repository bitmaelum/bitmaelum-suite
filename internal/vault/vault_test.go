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
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testingFilePath := "/my/dir/vault.json"

	fs = afero.NewMemMapFs()

	// inmemory vault
	v := New()
	assert.NotNil(t, v)

	// open/create new vault
	v, err := Create(testingFilePath, "secret")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	ok, _ := afero.Exists(fs, testingFilePath)
	assert.True(t, ok)

	// Reopen again
	v, err = Open(testingFilePath, "secret")
	assert.NoError(t, err)
	assert.NotNil(t, v)

	privKey, pubKey, _ := testing2.ReadTestKey("../../testdata/key-1.json")

	addr, _ := address.NewAddress("foobar!")
	acc := &AccountInfo{
		Address:  addr,
		Name:     "Foo Bar",
		Settings: nil,
		Keys: []KeyPair{
			{
				KeyPair: bmcrypto.KeyPair{
					Generator:   "",
					FingerPrint: pubKey.Fingerprint(),
					PrivKey:     *privKey,
					PubKey:      *pubKey,
				},
				Active: true,
			},
		},
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
	v, err = Open(testingFilePath, "incorrect password")
	assert.Errorf(t, err, "incorrect password")
	assert.Nil(t, v)

	// Open vault with correct password
	v, err = Open(testingFilePath, "secret")
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
	testingVault1FilePath := "/v1.json"
	v := NewPersistent(testingVault1FilePath, "foobar")
	assert.NotNil(t, v)

	assert.Equal(t, "foobar", v.password)
	assert.Equal(t, testingVault1FilePath, v.path)
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
	testingVault1FilePath := "/v1.json"
	testingVault2FilePath := "/v2.json"

	fs = afero.NewMemMapFs()

	// Cant open vaults
	v, err := Open(testingVault1FilePath, "foo")
	assert.Error(t, err)
	assert.Nil(t, v)
	v, err = Open(testingVault2FilePath, "bar")
	assert.Error(t, err)
	assert.Nil(t, v)

	// Create vaults
	_, err = Create(testingVault1FilePath, "foo")
	assert.NoError(t, err)
	_, err = Create(testingVault2FilePath, "bar")
	assert.NoError(t, err)

	// Open vaults
	v, err = Open(testingVault1FilePath, "foo")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	v, err = Open(testingVault2FilePath, "bar")
	assert.NoError(t, err)
	assert.NotNil(t, v)

	// Change password of vault
	v.SetPassword("anotherpass")
	err = v.Persist()
	assert.NoError(t, err)

	// Try open vault again
	v, err = Open(testingVault2FilePath, "bar")
	assert.Error(t, err)
	assert.Nil(t, v)
	v, err = Open(testingVault2FilePath, "anotherpass")
	assert.NoError(t, err)
	assert.NotNil(t, v)

}
