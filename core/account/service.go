package account

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Service is an account service that generates accounts on the REMOTE side (thus not client side)
type Service struct {
	repo Repository
}

// NewService initialises new service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateAccount creates new account for the given address and public key
func (s *Service) CreateAccount(addr address.HashAddress, pubKey string) error {
	if s.repo.Exists(addr) {
		return errors.New("account already exists")
	}

	err := s.repo.Create(addr)
	if err != nil {
		return err
	}

	_ = s.repo.CreateBox(addr, "inbox", "Inbox", "This is your regular inbox", 0)
	_ = s.repo.CreateBox(addr, "outbox", "Outbox", "All your outgoing messages will be stored here", 0)
	_ = s.repo.CreateBox(addr, "trash", "Trashcan", "Everything in here will be removed automatically after 30 days or when purged manually", 0)
	_ = s.repo.StorePubKey(addr, pubKey)

	return nil
}

// AccountExists checks if account exists for address
func (s *Service) AccountExists(addr address.HashAddress) bool {
	return s.repo.Exists(addr)
}

// GetPublicKeys retrieves the public keys for given address
func (s *Service) GetPublicKeys(addr address.HashAddress) []string {
	if !s.repo.Exists(addr) {
		return []string{}
	}

	pubKeys, err := s.repo.FetchPubKeys(addr)
	if err != nil {
		return []string{}
	}

	return pubKeys
}

// FetchMessageBoxes retrieves the message boxes based on the given query
func (s *Service) FetchMessageBoxes(addr address.HashAddress, query string) []core.MailBoxInfo {
	list, err := s.repo.FindBox(addr, query)
	if err != nil {
		return []core.MailBoxInfo{}
	}

	return list
}

// FetchListFromBox retrieves a list of message boxes
func (s *Service) FetchListFromBox(addr address.HashAddress, box string, offset int, limit int) []core.MessageList {
	list, err := s.repo.FetchListFromBox(addr, box, offset, limit)
	if err != nil {
		return []core.MessageList{}
	}

	return list
}

// GetFlags gets the flags for the given message
func (s *Service) GetFlags(addr address.HashAddress, box string, id string) ([]string, error) {
	return s.repo.GetFlags(addr, box, id)
}

// SetFlag sets a flag for a given message
func (s *Service) SetFlag(addr address.HashAddress, box string, id string, flag string) error {
	return s.repo.SetFlag(addr, box, id, flag)
}

// UnsetFlag unsets a flag for a given message
func (s *Service) UnsetFlag(addr address.HashAddress, box string, id string, flag string) error {
	return s.repo.UnsetFlag(addr, box, id, flag)
}
