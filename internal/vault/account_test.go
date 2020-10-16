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

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
)

func TestVaultAccount(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Store.Accounts, 0)

	acc := internal.AccountInfo{
		Address: "example!",
		Name:    "Example Account",
	}
	v.AddAccount(acc)
	assert.Len(t, v.Store.Accounts, 1)
	assert.Equal(t, "example!", v.Store.Accounts[0].Address)

	a, _ := address.NewAddress("example!")
	assert.True(t, v.HasAccount(*a))

	a, _ = address.NewAddress("notexists!")
	assert.False(t, v.HasAccount(*a))

	a, _ = address.NewAddress("example!")
	o, err := v.GetAccountInfo(*a)
	assert.NoError(t, err)
	assert.Equal(t, "example!", o.Address)

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

	acc := internal.AccountInfo{
		Default: false,
		Address: "acc1!",
	}
	v.AddAccount(acc)

	acc = internal.AccountInfo{
		Default: true,
		Address: "acc2!",
	}
	v.AddAccount(acc)

	acc = internal.AccountInfo{
		Default: false,
		Address: "acc3!",
	}
	v.AddAccount(acc)
	assert.Len(t, v.Store.Accounts, 3)

	da := v.GetDefaultAccount()
	assert.Equal(t, "acc2!", da.Address)

	// Without default set, pick the first
	v.Store.Accounts[1].Default = false
	da = v.GetDefaultAccount()
	assert.Equal(t, "acc1!", da.Address)

	v, _ = New("", []byte{})
	da = v.GetDefaultAccount()
	assert.Nil(t, da)
}
