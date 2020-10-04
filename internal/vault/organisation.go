package vault

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// AddOrganisation adds an organisation to the vault
func (v *Vault) AddOrganisation(organisation internal.OrganisationInfo) {
	v.Data.Organisations = append(v.Data.Organisations, organisation)
}

// GetOrganisationInfo tries to find the given organisation and returns the organisation from the vault
func (v *Vault) GetOrganisationInfo(org address.Hash) (*internal.OrganisationInfo, error) {

	for i := range v.Data.Organisations {
		h := address.NewHash(v.Data.Organisations[i].Addr)
		if h.String() == org.String() {
			return &v.Data.Organisations[i], nil
		}
	}

	return nil, errors.New("cannot find organisation")
}

// HasOrganisation returns true when the vault has an organisation for the given address
func (v *Vault) HasOrganisation(org address.Hash) bool {
	_, err := v.GetOrganisationInfo(org)

	return err == nil
}
