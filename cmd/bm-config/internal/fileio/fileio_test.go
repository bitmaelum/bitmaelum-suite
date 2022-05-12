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

package fileio

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSaveCertFiles(t *testing.T) {
	fs = afero.NewMemMapFs()

	var (
		err error
		ok  bool
	)

	config.Server.Server.CertFile = "/cert/cert.pem"
	config.Server.Server.KeyFile = "/cert/key.pem"

	err = SaveCertFiles("foo1", "bar1")
	assert.NoError(t, err)
	ok, _ = afero.Exists(fs, "/cert/cert.pem")
	assert.True(t, ok)
	ok, _ = afero.Exists(fs, "/cert/key.pem")
	assert.True(t, ok)

	err = SaveCertFiles("foo2", "bar2")
	assert.NoError(t, err)
	ok, _ = afero.Exists(fs, "/cert/cert.pem")
	assert.True(t, ok)
	ok, _ = afero.Exists(fs, "/cert/cert.pem.001")
	assert.True(t, ok)
	ok, _ = afero.Exists(fs, "/cert/key.pem")
	assert.True(t, ok)
	ok, _ = afero.Exists(fs, "/cert/key.pem.001")
	assert.True(t, ok)
}

func TestLoadSaveFile(t *testing.T) {
	fs = afero.NewMemMapFs()

	type MyStruct struct {
		A string
		B int
	}

	m := &MyStruct{
		A: "foo",
		B: 42,
	}
	err := SaveFile("/foo/bar.json", m)
	assert.NoError(t, err)

	ok, _ := afero.Exists(fs, "/foo/bar.json")
	assert.True(t, ok)

	m2 := &MyStruct{}
	err = LoadFile("/foo/bar.json", m2)
	assert.NoError(t, err)
	assert.Equal(t, "foo", m2.A)
	assert.Equal(t, 42, m2.B)
}
