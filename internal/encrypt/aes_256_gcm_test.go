package encrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func Test_EncryptDecryptJson(t *testing.T) {
	ts := &TestStruct{
		Foo: "foo",
		Bar: 42,
	}

	// MOck nonce generator for encrypt
	nonceGenerator = func(size int) ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil
	}

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

func Test_EncryptDecryptMessage(t *testing.T) {
	ts := "And now you do what they told ya"

	// MOck nonce generator for encrypt
	nonceGenerator = func(size int) ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil
	}

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

func Test_EncryptDecryptCatalog(t *testing.T) {
	type mockCatalog struct {
		Foo string
		Bar string
	}

	cat := mockCatalog{
		Foo: "foo",
		Bar: "bar",
	}

	// Mock key generator for encrypt
	nonceGenerator = func(size int) ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil
	}

	keyGenerator = func() ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}, nil
	}

	// With key and nonce, this should be the encrypted output
	dst := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd5, 0xc, 0xea, 0xf2, 0x5f, 0x49, 0x19, 0xcf, 0x9, 0xd2, 0x97, 0x7, 0xcf, 0xd8, 0x7f, 0x44, 0x7b, 0x7c, 0xc8, 0x5, 0x34, 0x14, 0xa3, 0x10, 0x28, 0xe3, 0xff, 0xa6, 0xe, 0x3a, 0xa9, 0x84, 0x10, 0x63, 0xa1, 0xd4, 0xd0, 0x2, 0x7f, 0x9c, 0x6a}

	// Encrypt
	key, encCatalog, err := CatalogEncrypt(cat)
	assert.Nil(t, err)
	dstKey, _ := keyGenerator()
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
