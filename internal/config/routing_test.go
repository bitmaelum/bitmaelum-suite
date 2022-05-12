// Copyright (c) 2022 BitMaelum Authors
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

package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	r, err := GenerateRouting()
	assert.NoError(t, err)

	assert.NotEmpty(t, r.RoutingID)
	assert.NotEmpty(t, r.KeyPair)
}

func TestReadSaveRouting(t *testing.T) {
	r, err := GenerateRouting()
	assert.NoError(t, err)

	fs = afero.NewMemMapFs()

	err = SaveRouting("/generated/routing.json", r)
	assert.NoError(t, err)

	err = ReadRouting("/generated/routing.json")
	assert.NoError(t, err)
	assert.Equal(t, r.RoutingID, Routing.RoutingID)
}
