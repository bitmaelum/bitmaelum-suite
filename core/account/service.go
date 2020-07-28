package account

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"time"
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

	_ = s.repo.CreateBox(addr, internal.BoxInbox)
	_ = s.repo.CreateBox(addr, internal.BoxOutbox)
	_ = s.repo.CreateBox(addr, internal.BoxTrash)
	_ = s.repo.StorePubKey(addr, pubKey)

	return nil
}

// AccountExists checks if account exists for address
func (s *Service) AccountExists(addr address.HashAddress) bool {
	return s.repo.Exists(addr)
}

// Deliver delivers a message (found in the processing queue) to the inbox of the given account
func (s *Service) Deliver(msgID string, addr address.HashAddress) error {
	return s.repo.SendToBox(addr, internal.BoxInbox, msgID)
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

// // FetchMessageBoxes retrieves the message boxes based on the given query
// func (s *Service) FetchMessageBoxes(addr address.HashAddress, query string) []message.MailBoxInfo {
// 	list, err := s.repo.FindBox(addr, query)
// 	if err != nil {
// 		return []message.MailBoxInfo{}
// 	}
//
// 	return list
// }

// FetchListFromBox retrieves a list of message boxes
func (s *Service) FetchListFromBox(addr address.HashAddress, box int, since time.Time, offset int, limit int) (*MessageList, error) {
	return s.repo.FetchListFromBox(addr, box, since, offset, limit)
}

// GetFlags gets the flags for the given message
func (s *Service) GetFlags(addr address.HashAddress, box int, id string) ([]string, error) {
	return s.repo.GetFlags(addr, box, id)
}

// SetFlag sets a flag for a given message
func (s *Service) SetFlag(addr address.HashAddress, box int, id string, flag string) error {
	return s.repo.SetFlag(addr, box, id, flag)
}

// UnsetFlag unsets a flag for a given message
func (s *Service) UnsetFlag(addr address.HashAddress, box int, id string, flag string) error {
	return s.repo.UnsetFlag(addr, box, id, flag)
}

type BoxInfo struct {
	ID    int   `json:"id"`
	Total int   `json:"total"`
}

func (s *Service) GetAllBoxes(addr address.HashAddress) ([]BoxInfo, error) {
	return s.repo.GetAllBoxes(addr)
}

func (s *Service) ExistsBox(addr address.HashAddress, box int) bool {
	return s.repo.ExistsBox(addr, box)
}
