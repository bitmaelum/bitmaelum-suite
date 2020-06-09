package main

import (
    "context"
    "crypto"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "github.com/sirupsen/logrus"
    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
    "io/ioutil"
    "net/http"
    "os"
    "time"
)

const StagingAcmeDir = "https://acme-staging-v02.api.letsencrypt.org/directory"

func main() {
    client := &acme.Client{
        DirectoryURL: StagingAcmeDir,
        Key: rs256
    }

    m := &autocert.Manager{
        Client:     client,
        Cache:      autocert.DirCache(".certs"),
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist("mailv2.ngrok.io"),
    }
    s := &http.Server{
        Addr:      ":2425",
        TLSConfig: m.TLSConfig(),
    }
    go func() {
        s.ListenAndServeTLS("", "")
    }()

    defer s.Close()


    ctx, _ := context.WithTimeout(context.Background(), 60 * time.Second)

    client.Key
    account, err := client.Register(ctx, &acme.Account{Contact: []string{"jthijssen@noxlogic.nl"}}, acme.AcceptTOS)
    if err != nil {
       panic(err)
    }
    fmt.Printf("%#v", account)

    domain := "mailv2.ngrok.io"

    certKey, err := loadOrGenerateKey()
    if err != nil {
        panic(err)
    }


    certRequest := &x509.CertificateRequest{
        Subject:  pkix.Name{CommonName: domain},
        DNSNames: []string{domain},
    }

    csr, err := x509.CreateCertificateRequest(rand.Reader, certRequest, certKey)
    if err != nil {
        panic(err)
    }

    der, url, err := client.CreateCert(ctx, csr, 90 * 24 * time.Hour, true)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v", der)
    fmt.Printf("%#v", url)
}

func generateKey() (*ecdsa.PrivateKey, error) {
    key, err := LoadKey(file)
    csdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        logrus.Error(err)
        return nil, err
    }

    bytes, err := x509.MarshalECPrivateKey(csdsaKey)
    if err != nil {
        return nil, err
    }
    f, err := os.OpenFile("acme.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return nil, err
    }
    b := &pem.Block{Type: "EC PRIVATE KEY", Bytes: bytes}
    if err := pem.Encode(f, b); err != nil {
        f.Close()
        return nil, err
    }
    _ = f.Close()

    return csdsaKey, nil
}

func LoadKey(file string) (crypto.Signer, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		logrus.Infof("open %s: no such file or directory, will re-generate key file", file)
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil {
		logrus.Errorf("failed to decode pem file: no key found")
		return nil, errors.Errorf("unsupported type: %s", block.Type)
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(block.Bytes)
	default:
		logrus.Errorf("unsupported type: %s", block.Type)
		return nil, errors.Errorf("unsupported type: %s", block.Type)
	}
}
