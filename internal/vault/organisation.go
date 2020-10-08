package vault

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// AddOrganisation adds an organisation to the vault
func (v *Vault) AddOrganisation(organisation internal.OrganisationInfo) {
	v.Store.Organisations = append(v.Store.Organisations, organisation)
}

// GetOrganisationInfo tries to find the given organisation and returns the organisation from the vault
func (v *Vault) GetOrganisationInfo(org hash.Hash) (*internal.OrganisationInfo, error) {

	for i := range v.Store.Organisations {
		h := hash.New(v.Store.Organisations[i].Addr)
		if h.String() == org.String() {
			return &v.Store.Organisations[i], nil
		}
	}

	return nil, errors.New("cannot find organisation")
}

// HasOrganisation returns true when the vault has an organisation for the given address
func (v *Vault) HasOrganisation(org hash.Hash) bool {
	_, err := v.GetOrganisationInfo(org)

	return err == nil
}
