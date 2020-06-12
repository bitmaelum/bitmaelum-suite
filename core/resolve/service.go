package resolve

import (
    "crypto/sha256"
    "encoding/hex"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/encrypt"
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

    // @TODO: We maybe should sign with a different algo? Otherwise we use the same one for all systems
    privKey, err := encrypt.PEMToPrivKey([]byte(acc.PrivKey))
    if err != nil {
        return err
    }

    // Sign resolve address
    hash := sha256.Sum256([]byte(resolveAddress))
    signature, err := encrypt.Sign(privKey, hash[:])
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

