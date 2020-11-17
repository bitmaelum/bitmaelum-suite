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

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

func TestVaultAccount(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Store.Accounts, 0)

	addr, _ := address.NewAddress("example!")
	acc := AccountInfo{
		Address: *addr,
		Name:    "Example Account",
	}
	v.AddAccount(acc)
	assert.Len(t, v.Store.Accounts, 1)
	assert.Equal(t, "example!", v.Store.Accounts[0].Address.String())

	a, _ := address.NewAddress("example!")
	assert.True(t, v.HasAccount(*a))

	a, _ = address.NewAddress("notexists!")
	assert.False(t, v.HasAccount(*a))

	a, _ = address.NewAddress("example!")
	o, err := v.GetAccountInfo(*a)
	assert.NoError(t, err)
	assert.Equal(t, "example!", o.Address.String())

	a, _ = address.NewAddress("notexist!")
	o, err = v.GetAccountInfo(*a)
	assert.Error(t, err)
	assert.Nil(t, o)

	assert.Len(t, v.Store.Accounts, 1)
	a, _ = address.NewAddress("notexist!")
	v.RemoveAccount(*a)
	assert.Len(t, v.Store.Accounts, 1)

	a, _ = address.NewAddress("example!")
	v.RemoveAccount(*a)
	assert.Len(t, v.Store.Accounts, 0)
}

func TestVaultGetDefaultAccount(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Store.Accounts, 0)

	addr, _ := address.NewAddress("acc1!")
	acc := AccountInfo{
		Default: false,
		Address: *addr,
	}
	v.AddAccount(acc)

	addr, _ = address.NewAddress("acc2!")
	acc = AccountInfo{
		Default: true,
		Address: *addr,
	}
	v.AddAccount(acc)

	addr, _ = address.NewAddress("acc3!")
	acc = AccountInfo{
		Default: false,
		Address: *addr,
	}
	v.AddAccount(acc)
	assert.Len(t, v.Store.Accounts, 3)

	da := v.GetDefaultAccount()
	assert.Equal(t, "acc2!", da.Address.String())

	// Without default set, pick the first
	v.Store.Accounts[1].Default = false
	da = v.GetDefaultAccount()
	assert.Equal(t, "acc1!", da.Address.String())

	v, _ = New("", []byte{})
	da = v.GetDefaultAccount()
	assert.Nil(t, da)
}

func TestInfoToOrg(t *testing.T) {
	info := &OrganisationInfo{
		Addr:        "foo",
		FullName:    "bar",
		PrivKey:     bmcrypto.PrivKey{},
		PubKey:      bmcrypto.PubKey{},
		Pow:         proofofwork.New(22, "foobar", 1234),
		Validations: nil,
	}

	assert.Equal(t, "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae", info.ToOrg().Hash.String())
}

func TestAccountInfoAddressHash(t *testing.T) {
	addr, _ := address.NewAddress("example!")
	info := &AccountInfo{
		Default: false,
		Address: *addr,
		Name:    "John DOe",
	}

	assert.Equal(t, "2244643da7475120bf84d744435d15ea297c36ca165ea0baaa69ec818d0e952f", info.Address.Hash().String())
}
