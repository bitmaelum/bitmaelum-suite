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

package message

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
)

var src = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sagittis sapien non risus sollicitudin malesuada."
var dst = []byte{0x78, 0xda, 0x14, 0xc6, 0xd1, 0x89, 0x3, 0x31, 0xc, 0x4, 0xd0, 0x56, 0xa6, 0x80, 0x63, 0x2b, 0xb9, 0x14, 0x21, 0x6c, 0xb1, 0xc, 0xd8, 0x92, 0xf1, 0x48, 0xfd, 0x87, 0xfc, 0xbd, 0xff, 0xbc, 0xbe, 0xc1, 0xa3, 0xde, 0x98, 0xb9, 0xf2, 0x42, 0x2c, 0xd8, 0xf6, 0xfa, 0xc3, 0xc8, 0x90, 0x8f, 0xf2, 0xea, 0xb, 0x9b, 0x3c, 0xd4, 0x60, 0xbc, 0xf0, 0xc5, 0x7a, 0xf0, 0x31, 0x1f, 0x1e, 0x26, 0xc8, 0x5e, 0x56, 0xf1, 0x87, 0x43, 0xf, 0x44, 0x6, 0x2e, 0xd5, 0x82, 0x72, 0x2d, 0xe, 0x56, 0x4f, 0x6, 0xb6, 0x2d, 0x57, 0xdb, 0xb4, 0xe7, 0x1b, 0x0, 0x0, 0xff, 0xff, 0xbf, 0x8d, 0x2b, 0x4b}

var src1 = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
var dst1 = []byte{0x78, 0xda, 0x72, 0x1c, 0xe4, 0x0, 0x10, 0x0, 0x0, 0xff, 0xff, 0xc7, 0xa4, 0x28, 0xa1}

func TestCompress(t *testing.T) {
	r := bytes.NewBufferString(src)
	compressedBytes, _ := ioutil.ReadAll(ZlibCompress(r))
	assert.Equal(t, 0, bytes.Compare(compressedBytes, dst))

	r = bytes.NewBufferString(src1)
	compressedBytes, _ = ioutil.ReadAll(ZlibCompress(r))
	assert.Equal(t, 0, bytes.Compare(compressedBytes, dst1))
}

func TestDecompress(t *testing.T) {
	r, err := ZlibDecompress(bytes.NewReader(dst))
	assert.NoError(t, err)
	decompressedBytes, _ := ioutil.ReadAll(r)
	assert.Equal(t, 0, bytes.Compare(decompressedBytes, []byte(src)))

	r, _ = ZlibDecompress(bytes.NewReader(dst1))
	decompressedBytes, _ = ioutil.ReadAll(r)
	assert.Equal(t, 0, bytes.Compare(decompressedBytes, []byte(src1)))
}

func TestEncryptAndCompress(t *testing.T) {
	iv, key, _ := bmcrypto.GenerateIvAndKey()

	// Open our attachment PNG
	f, err := os.Open("../../testdata/attachment.png")
	assert.NoError(t, err)

	// Compress then encrypt
	r := ZlibCompress(f)
	r, _ = bmcrypto.GetAesEncryptorReader(iv, key, r)

	encData, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, len(encData), 5407)

	// Decrypt, then decompress
	r, err = bmcrypto.GetAesDecryptorReader(iv, key, bytes.NewReader(encData))
	assert.NoError(t, err)
	r, err = ZlibDecompress(r)
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, data[:4], []byte{0x89, 'P', 'N', 'G'})
}
