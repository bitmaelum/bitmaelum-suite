package encrypt

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/message"
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
	data, err := EncryptJson(key, ts)
	assert.Nil(t, err)
	assert.Equal(t, dst, data)

	// Decrypt again
	ts1 := &TestStruct{}
	err = DecryptJson(key, data, &ts1)
	assert.Nil(t, err)
	assert.Equal(t, "foo", ts1.Foo)
	assert.Equal(t, 42, ts1.Bar)

	// Validate that another key does not decrypt
	key[0] ^= 80
	ts2 := &TestStruct{}
	err = DecryptJson(key, data, &ts2)
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
	data, err := EncryptMessage(key, []byte(ts))
	assert.Nil(t, err)
	assert.Equal(t, dst, data)

	// Decrypt again
	plainText, err := DecryptMessage(key, data)
	assert.Nil(t, err)
	assert.Equal(t, ts, string(plainText))

	// Validate that another key does not decrypt
	key[0] ^= 80
	_, err = DecryptMessage(key, data)
	assert.EqualError(t, err, "cipher: message authentication failed")
}

func Test_EncryptDecryptCatalog(t *testing.T) {
	cat := &message.Catalog{
		From: struct {
			Address      string           `json:"address"`
			Name         string           `json:"name"`
			Organisation string           `json:"organisation"`
			ProofOfWork  core.ProofOfWork `json:"proof_of_work"`
			PublicKey    string           `json:"public_key"`
		}{
			Address:      "bitmaelum!",
			Name:         "Test user",
			Organisation: "BitMaelum",
		},
		ThreadId: "1234",
		Subject:  "Our subject matters",
	}

	// Mock key generator for encrypt
	nonceGenerator = func(size int) ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil
	}

	keyGenerator = func() ([]byte, error) {
		return []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}, nil
	}

	// With key and nonce, this should be the encrypted output
	dst := []byte("AQIDBAUGBwgJCgsM1QzK718GAdcUn5lBh4hYVnp8yAU0HKVfNHIZPnuBXpCn2VjaqP+qPskB9XcEeRWGqyG/ycW/8NMziFYWXd3SJNSeUrz6tsjKO2ASAlsqkiqaNg9g1CghOkR3mt4xSVBN0EuavGi5UpuN5abBPr8cwtom5JfywudPkKn0o2o6teX0cY9L5pCECa9spdgCRme23lRHEXTqHDMeDI+LlT8gfpInuCkeSxb00to/7cKAweQFDd/cUVtsWmmI+y8BllzosMz2cmiwasmOscrWoM102RnjHQjHo8F6SbEUkKUSIKXZGYI4pRT5hrgR8iuAXzpdaGCI7/2otCIQZJoHu8kPHVKIOQQXdX71AqD/WOX8tJz8DJ2xvts1/UweOv0C1j7AxB+XbVByXHsxVu69LphNPy8PASQ81o4sYqti")

	// Encrypt
	key, encCatalog, err := EncryptCatalog(*cat)
	assert.Nil(t, err)
	dstKey, _ := keyGenerator()
	assert.Equal(t, dstKey, key)
	assert.Equal(t, dst, encCatalog)

	// Decrypt again
	cat2, err := DecryptCatalog(key, encCatalog)
	assert.Nil(t, err)
	assert.Equal(t, "1234", cat2.ThreadId)

	// Validate that another key does not decrypt
	key[0] ^= 80
	_, err = DecryptCatalog(key, encCatalog)
	assert.EqualError(t, err, "cipher: message authentication failed")
}
