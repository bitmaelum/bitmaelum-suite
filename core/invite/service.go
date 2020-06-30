package invite

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"time"
)

// Service is the invitation service
type Service struct {
	repo Repository
}

// NewInviteService create new invitation service based on the given repository
func NewInviteService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateInvite creates a new invitation in the repository
func (s *Service) CreateInvite(addr core.HashAddress, expiry time.Duration) (string, error) {
	return s.repo.CreateInvite(addr, expiry)
}

// GetInvite retrieves an invitation from the repository
func (s *Service) GetInvite(addr core.HashAddress) (string, error) {
	return s.repo.GetInvite(addr)
}

// RemoveInvite deletes an invitation from the repository
func (s *Service) RemoveInvite(addr core.HashAddress) error {
	return s.repo.RemoveInvite(addr)
}
