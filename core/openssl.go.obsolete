package core

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
)

func CreateNewKeyPair(bits int) (*rsa.PrivateKey, error) {
    key, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, err
    }

    return key, err
}

func LoadPubKey(path string) (*rsa.PublicKey, error) {
    derData, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode([]byte(derData))
    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return key.(*rsa.PublicKey), err
}

func LoadPrivKey(path string) (*rsa.PrivateKey, error) {
    derData, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode([]byte(derData))
    return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func SavePrivKey(path string, privKey *rsa.PrivateKey) error {
    pemData := x509.MarshalPKCS1PrivateKey(privKey)
    pemDataBytes := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: pemData,
    })

    err := ioutil.WriteFile(path, pemDataBytes, 0600)
    return err;
}
