package account

import (
    "errors"
)

type Service struct {
    repo Repository
}

func NewAccountService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (s *Service) CreateAccount(hash string, pubKey string) error {
    if s.repo.Exists(hash) {
        return errors.New("account already exists")
    }

    err := s.repo.Create(hash)
    if err != nil {
        return err
    }

    _ = s.repo.CreateBox(hash, "inbox", "This is your regular inbox", 0)
    _ = s.repo.CreateBox(hash, "outbox", "All your outgoing emails will be stored here", 0)
    _ = s.repo.CreateBox(hash, "trash", "Trashcan. Everything in here will be removed automatically after 30 days or when purged manually", 0)
    _ = s.repo.StorePubKey(hash, ".pubkey", []byte(pubKey))

    return nil
}
