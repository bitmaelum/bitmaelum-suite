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
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// AddOrganisation adds an organisation to the vault
func (v *Vault) AddOrganisation(organisation internal.OrganisationInfo) {
	v.Store.Organisations = append(v.Store.Organisations, organisation)
}

// GetOrganisationInfo tries to find the given organisation and returns the organisation from the vault
func (v *Vault) GetOrganisationInfo(orgHash hash.Hash) (*internal.OrganisationInfo, error) {

	for i := range v.Store.Organisations {
		h := hash.New(v.Store.Organisations[i].Addr)
		if h.String() == orgHash.String() {
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
