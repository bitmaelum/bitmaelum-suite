package vault

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
)

func TestVault_AddOrganisation(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Data.Organisations, 0)

	org := internal.OrganisationInfo{
		Addr: "example",
		Name: "Example Org",
	}
	v.AddOrganisation(org)
	assert.Len(t, v.Data.Organisations, 1)
	assert.Equal(t, "example", v.Data.Organisations[0].Addr)

	a, _ := address.NewOrganisationHash("example")
	assert.True(t, v.HasOrganisation(*a))

	a, _ = address.NewOrganisationHash("notexist")
	assert.False(t, v.HasOrganisation(*a))

	a, _ = address.NewOrganisationHash("example")
	o, err := v.GetOrganisationInfo(*a)
	assert.NoError(t, err)
	assert.Equal(t, "example", o.Addr)

	a, _ = address.NewOrganisationHash("notexist")
	o, err = v.GetOrganisationInfo(*a)
	assert.Error(t, err)
	assert.Nil(t, o)
}
