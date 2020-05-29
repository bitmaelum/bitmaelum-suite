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
    "github.com/jaytaph/mailv2/core/account"
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

func (s *Service) GetInfo(addr core.Address) (*ResolveInfo, error) {
    return s.repo.Retrieve(addr.ToHash())
}

func (s *Service) UploadInfo(acc account.Account, resolveAddress string) error {
    block, _ := pem.Decode([]byte(acc.PrivKey))
    privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return err
    }

    hash := sha256.Sum256([]byte(resolveAddress))
    signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
    if err != nil {
        return err
    }

    return s.repo.Upload(
        core.StringToHash(acc.Address),
        acc.PubKey,
        resolveAddress,
        hex.EncodeToString(signature),
    )
}

