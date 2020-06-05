package resolve

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/hex"
    "encoding/pem"
    "github.com/jaytaph/mailv2/core"
)

type Service struct {
    repo Repository
}

type ResolveInfo struct {
    Hash        string `json:"hash"`
    PublicKey   string `json:"public_key"`
    Address     string `json:"address"`
}

func KeyRetrievalService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

// Resolve an address
func (s *Service) Resolve(addr core.Address) (*ResolveInfo, error) {
    return s.repo.Resolve(addr.Hash())
}

// Upload resolve information to a service
func (s *Service) UploadInfo(acc core.AccountInfo, resolveAddress string) error {
    // Get private key
    block, _ := pem.Decode([]byte(acc.PrivKey))
    privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return err
    }

    // Sign resolve address
    hash := sha256.Sum256([]byte(resolveAddress))
    signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
    if err != nil {
        return err
    }

    // And upload
    return s.repo.Upload(
        core.StringToHash(acc.Address),
        acc.PubKey,
        resolveAddress,
        hex.EncodeToString(signature),
    )
}

