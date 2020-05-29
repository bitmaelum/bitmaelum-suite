package account

import (
    "errors"
    "github.com/jaytaph/mailv2/core"
)

type Service struct {
    repo Repository
}

// Create new service
func AccountService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

// Create new account for the given address and public key
func (s *Service) CreateAccount(addr core.HashAddress, pubKey string) error {
    if s.repo.Exists(addr) {
        return errors.New("account already exists")
    }

    err := s.repo.Create(addr)
    if err != nil {
        return err
    }

    _ = s.repo.CreateBox(addr, "inbox", "This is your regular inbox", 0)
    _ = s.repo.CreateBox(addr, "outbox", "All your outgoing messages will be stored here", 0)
    _ = s.repo.CreateBox(addr, "trash", "Trashcan. Everything in here will be removed automatically after 30 days or when purged manually", 0)
    _ = s.repo.StorePubKey(addr, []byte(pubKey))

    return nil
}

// Check if account exists for address
func (s *Service) AccountExists(addr core.HashAddress) bool {
    return s.repo.Exists(addr)
}

// Retrieve the public key for given address
func (s *Service) GetPublicKey(addr core.HashAddress) string {
    if ! s.repo.Exists(addr) {
        return ""
    }

    pubKey, err := s.repo.FetchPubKey(addr)
    if err != nil {
        return ""
    }

    return string(pubKey)
}

