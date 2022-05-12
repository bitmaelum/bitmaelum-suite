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

package bmcrypto

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func mockGenerator(size int) ([]byte, error) {
	fixedBytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	return fixedBytes[:size], nil
}

func TestEncryptDecryptJson(t *testing.T) {
	ts := &TestStruct{
		Foo: "foo",
		Bar: 42,
	}

	// Mock nonce generator for encrypt
	keyGenerator = mockGenerator

	// Key we are using to encrypt
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	// Mock nonce
	// With key and nonce, this should be the encrypted output
	dst := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd5, 0xc, 0xca, 0xf2, 0x5f, 0x49, 0x19, 0xcf, 0x9, 0xd2, 0x97, 0x7, 0xcf, 0xd8, 0x5f, 0x44, 0x7b, 0x7c, 0xc8, 0x13, 0x64, 0x8, 0x6e, 0xf2, 0x88, 0x62, 0xb9, 0x25, 0x35, 0x67, 0xe5, 0x63, 0x89, 0xd1, 0x9e, 0x85, 0xd9, 0x6d}

	// Encrypt
	data, err := JSONEncrypt(key, ts)
	assert.Nil(t, err)
	assert.Equal(t, dst, data)

	// Decrypt again
	ts1 := &TestStruct{}
	err = JSONDecrypt(key, data, &ts1)
	assert.Nil(t, err)
	assert.Equal(t, "foo", ts1.Foo)
	assert.Equal(t, 42, ts1.Bar)

	// Validate that another key does not decrypt
	key[0] ^= 80
	ts2 := &TestStruct{}
	err = JSONDecrypt(key, data, &ts2)
	assert.EqualError(t, err, "cipher: message authentication failed")
}

func TestEncryptDecryptMessage(t *testing.T) {
	ts := "And now you do what they told ya"

	// MOck nonce generator for encrypt
	keyGenerator = mockGenerator

	// Key we are using to encrypt
	key := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

	// With key and nonce, this should be the encrypted output
	dst := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xef, 0x40, 0xc8, 0xbd, 0x5e, 0x4, 0x54, 0xcd, 0x16, 0xd2, 0x8d, 0x5, 0x87, 0x95, 0x1d, 0x52, 0x61, 0x3f, 0x86, 0x7, 0x22, 0x1d, 0xb4, 0x4b, 0x75, 0x63, 0x1a, 0x27, 0x72, 0x80, 0x5, 0xdd, 0x8a, 0x16, 0x72, 0x8c, 0x2d, 0x98, 0xfa, 0xf6, 0xe5, 0xc1, 0x6, 0x7a, 0xdd, 0x42, 0xae, 0x90}

	// Encrypt
	data, err := MessageEncrypt(key, []byte(ts))
	assert.Nil(t, err)
	assert.Equal(t, dst, data)

	// Decrypt again
	plainText, err := MessageDecrypt(key, data)
	assert.Nil(t, err)
	assert.Equal(t, ts, string(plainText))

	// Validate that another key does not decrypt
	key[0] ^= 80
	_, err = MessageDecrypt(key, data)
	assert.EqualError(t, err, "cipher: message authentication failed")
}

func TestEncryptDecryptCatalog(t *testing.T) {
	type mockCatalog struct {
		Foo string
		Bar string
	}

	cat := mockCatalog{
		Foo: "foo",
		Bar: "bar",
	}

	// // Mock key generator for encrypt
	keyGenerator = mockGenerator

	// With key and nonce, this should be the encrypted output
	dst := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd5, 0xc, 0xea, 0xf2, 0x5f, 0x49, 0x19, 0xcf, 0x9, 0xd2, 0x97, 0x7, 0xcf, 0xd8, 0x7f, 0x44, 0x7b, 0x7c, 0xc8, 0x5, 0x34, 0x14, 0xa3, 0x10, 0x28, 0xe3, 0xff, 0xa6, 0xe, 0x3a, 0xa9, 0x84, 0x10, 0x63, 0xa1, 0xd4, 0xd0, 0x2, 0x7f, 0x9c, 0x6a}

	// Encrypt
	key, encCatalog, err := CatalogEncrypt(cat)
	assert.Nil(t, err)
	dstKey, _ := keyGenerator(32)
	assert.Equal(t, dstKey, key)
	assert.Equal(t, dst, encCatalog)

	// Decrypt again
	cat2 := &mockCatalog{}
	err = CatalogDecrypt(key, encCatalog, cat2)
	assert.Nil(t, err)
	assert.Equal(t, "bar", cat2.Bar)

	// Validate that another key does not decrypt
	key[0] ^= 80
	err = CatalogDecrypt(key, encCatalog, cat2)
	assert.EqualError(t, err, "cipher: message authentication failed")
}

func TestRandomGenerator(t *testing.T) {
	b1, err := randomKeyGenerator(32)
	assert.NoError(t, err)
	assert.Len(t, b1, 32)

	b2, err := randomKeyGenerator(32)
	assert.NoError(t, err)
	assert.Len(t, b2, 32)

	assert.NotEqual(t, b1, b2)
}

func TestEncryptor(t *testing.T) {
	iv := []byte{0xa4, 0xc0, 0x44, 0xc6, 0x3c, 0xdd, 0x9c, 0xe6, 0xae, 0x62, 0xd7, 0xf3, 0x84, 0x6a, 0x21, 0x2e}
	key := []byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7, 0x49, 0xc2, 0x7, 0x87, 0xcf, 0xfb, 0xac, 0x10, 0xb0, 0xc4, 0xcc, 0xf0, 0xee, 0xb3, 0xe9, 0xb8, 0x9e, 0x7c, 0xfb, 0xce, 0x58, 0xc2, 0x2f}

	r1 := bytes.NewBufferString("foo bar baz")
	r2, err := GetAesEncryptorReader(iv, key, r1)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r2)
	assert.Equal(t, []byte{0xc3, 0x7, 0xbd, 0xa2, 0x65, 0x50, 0x6, 0x90, 0xa3, 0x39, 0x81}, b)
	assert.NoError(t, err)

	r3 := bytes.NewReader(b)
	r4, err := GetAesDecryptorReader(iv, key, r3)
	assert.NoError(t, err)

	b, err = ioutil.ReadAll(r4)
	assert.NoError(t, err)
	assert.Equal(t, "foo bar baz", string(b))
}

func TestGenerateIvAndKey(t *testing.T) {
	iv, key, err := GenerateIvAndKey()
	assert.NoError(t, err)
	assert.Len(t, iv, 16)
	assert.Len(t, key, 32)
}

func TestKeySizes(t *testing.T) {
	iv := []byte{0xa4, 0xc0, 0x44, 0xc6, 0x3c, 0xdd, 0x9c, 0xe6, 0xae, 0x62, 0xd7, 0xf3, 0x84, 0x6a, 0x21, 0x2e}
	key := []byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7, 0x49, 0xc2, 0x7, 0x87, 0xcf, 0xfb, 0xac, 0x10, 0xb0, 0xc4, 0xcc, 0xf0, 0xee, 0xb3, 0xe9, 0xb8, 0x9e, 0x7c, 0xfb, 0xce, 0x58, 0xc2, 0x2f}

	r1 := bytes.NewBufferString("foo bar baz")
	_, err := GetAesEncryptorReader(iv, []byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7}, r1)
	assert.Error(t, err)

	r1 = bytes.NewBufferString("foo bar baz")
	_, err = GetAesEncryptorReader([]byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7}, key, r1)
	assert.Error(t, err)

	r1 = bytes.NewBufferString("foo bar baz")
	_, err = GetAesDecryptorReader(iv, []byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7}, r1)
	assert.Error(t, err)

	r1 = bytes.NewBufferString("foo bar baz")
	_, err = GetAesDecryptorReader([]byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7}, key, r1)
	assert.Error(t, err)
}

func TestCreateCatalogKey(t *testing.T) {
	// MOck nonce generator for encrypt
	keyGenerator = mockGenerator

	b, err := CreateCatalogKey()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}, b)
}

func TestJSONEncryptDecrypt(t *testing.T) {
	key := []byte{0xf2, 0x3f, 0x97, 0x75, 0xc6, 0x40, 0xfa, 0x3e, 0xf7, 0x49, 0xc2, 0x7, 0x87, 0xcf, 0xfb, 0xac, 0x10, 0xb0, 0xc4, 0xcc, 0xf0, 0xee, 0xb3, 0xe9, 0xb8, 0x9e, 0x7c, 0xfb, 0xce, 0x58, 0xc2, 0x2f}

	b, err := JSONEncrypt(key, make(chan int))
	assert.Error(t, err)
	assert.Nil(t, b)

	var v string
	err = JSONDecrypt(key, []byte("foobar"), &v)
	assert.Error(t, err)
}
