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
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestVaultAddOrganisation(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)

	assert.Len(t, v.Store.Organisations, 0)

	org := internal.OrganisationInfo{
		Addr:     "example",
		FullName: "Example Org",
	}
	v.AddOrganisation(org)
	assert.Len(t, v.Store.Organisations, 1)
	assert.Equal(t, "example", v.Store.Organisations[0].Addr)

	a := hash.New("example")
	assert.True(t, v.HasOrganisation(a))

	a = hash.New("notexist")
	assert.False(t, v.HasOrganisation(a))

	a = hash.New("example")
	o, err := v.GetOrganisationInfo(a)
	assert.NoError(t, err)
	assert.Equal(t, "example", o.Addr)

	a = hash.New("notexist")
	o, err = v.GetOrganisationInfo(a)
	assert.Error(t, err)
	assert.Nil(t, o)
}
