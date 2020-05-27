package keys

import (
    "github.com/jaytaph/mailv2/core"
)

type Service struct {
    repo Repository
}

func KeyRetrievalService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (s *Service) GetPublicKey(email string) (string, error) {
    return s.repo.Retrieve(core.HashEmail(email))
}

