package vault

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
)

func TestVault_Account(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Data.Accounts, 0)

	acc := internal.AccountInfo{
		Address: "example!",
		Name:    "Example Account",
	}
	v.AddAccount(acc)
	assert.Len(t, v.Data.Accounts, 1)
	assert.Equal(t, "example!", v.Data.Accounts[0].Address)

	a, _ := address.New("example!")
	assert.True(t, v.HasAccount(*a))

	a, _ = address.New("notexists!")
	assert.False(t, v.HasAccount(*a))

	a, _ = address.New("example!")
	o, err := v.GetAccountInfo(*a)
	assert.NoError(t, err)
	assert.Equal(t, "example!", o.Address)

	a, _ = address.New("notexist!")
	o, err = v.GetAccountInfo(*a)
	assert.Error(t, err)
	assert.Nil(t, o)

	assert.Len(t, v.Data.Accounts, 1)
	a, _ = address.New("notexist!")
	v.RemoveAccount(*a)
	assert.Len(t, v.Data.Accounts, 1)

	a, _ = address.New("example!")
	v.RemoveAccount(*a)
	assert.Len(t, v.Data.Accounts, 0)
}

func TestVault_GetDefaultAccount(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Data.Accounts, 0)

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
	assert.Len(t, v.Data.Accounts, 3)

	da := v.GetDefaultAccount()
	assert.Equal(t, "acc2!", da.Address)

	// Without default set, pick the first
	v.Data.Accounts[1].Default = false
	da = v.GetDefaultAccount()
	assert.Equal(t, "acc1!", da.Address)

	v, _ = New("", []byte{})
	da = v.GetDefaultAccount()
	assert.Nil(t, da)
}
