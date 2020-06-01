package encrypt

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/json"
    "github.com/jaytaph/mailv2/core/catalog"
    "github.com/jaytaph/mailv2/core/encode"
    "io"
)

// Encrypt json data with AES256
func EncryptJson(key, iv []byte, data interface{}) ([]byte, error) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return nil, err
    }

    plaintext := []byte(jsonData)
    cfb := cipher.NewCFBEncrypter(block, iv)
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)

    return ciphertext, nil
}

// Decrypt AES256 data back into json data
func DecryptJson(key, iv []byte, ciphertext []byte, v interface{}) error {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return err
    }

    cfb := cipher.NewCFBDecrypter(block, iv)
    plaintext := make([]byte, len(ciphertext))
    cfb.XORKeyStream(plaintext, ciphertext)

    return json.Unmarshal(plaintext, &v)
}

func EncryptData(key, iv []byte, r io.Reader, w *io.Writer) {
}


// Encrypt a catalog data. Use random key and iv.
func EncryptCatalog(catalog catalog.Catalog) ([]byte, []byte, []byte, error) {
    var err error

    catalogKey := make([]byte, 32)
    _, err = rand.Read(catalogKey)
    if err != nil {
        return nil, nil, nil, err
    }
    catalogIv := make([]byte, 16)
    _, err = rand.Read(catalogIv)
    if err != nil {
        return nil, nil, nil, err
    }

    ciphertext, err := EncryptJson(catalogKey, catalogIv, catalog)

    return catalogKey, catalogIv, encode.Encode(ciphertext), nil
}
