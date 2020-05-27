package encrypt

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
)

func EncryptKey(pubKey, data []byte) ([]byte, error) {
    block, _ := pem.Decode(pubKey)
    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    return rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), data)
}
