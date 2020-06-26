package invite

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"time"
)

type Service struct {
	repo Repository
}

// Create new service
func NewInviteService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateInvite(addr core.HashAddress, expiry time.Duration) (string, error) {
	return s.repo.CreateInvite(addr, expiry)
}

func (s *Service) GetInvite(addr core.HashAddress) (string, error) {
	return s.repo.GetInvite(addr)
}

func (s *Service) RemoveInvite(addr core.HashAddress) error {
	return s.repo.RemoveInvite(addr)
}
