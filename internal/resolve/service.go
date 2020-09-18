package resolve

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
)

// Service represents a resolver service tied to a specific repository
type Service struct {
	repo Repository
}

// Info is a structure returned by the external resolver system
type Info struct {
	Hash      string          `json:"hash"`       // Hash of the email address
	PublicKey bmcrypto.PubKey `json:"public_key"` // PublicKey of the user
	Routing   string          `json:"routing"`    // Server where this email address resides
	Pow       string          `json:"pow"`        // Proof of work
}

// KeyRetrievalService initialises a key retrieval service.
func KeyRetrievalService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Resolve resolves an address.
func (s *Service) Resolve(addr address.HashAddress) (*Info, error) {
	logrus.Debugf("Resolving %s", addr.String())
	info, err := s.repo.Resolve(addr)
	if err != nil {
		logrus.Debugf("Error while resolving %s: %s", addr.String(), err)
	}

	return info, err
}

// UploadInfo uploads resolve information to one (or more) resolvers
func (s *Service) UploadInfo(info internal.AccountInfo) error {
	hashAddr, err := address.NewHash(info.Address)
	if err != nil {
		return err
	}

	return s.repo.Upload(&Info{
		Hash:      hashAddr.String(),
		PublicKey: info.PubKey,
		Routing:   info.Routing,
		Pow:       info.Pow.String(),
	}, info.PrivKey, info.Pow)
}

// generateSignature generates a signature with the accounts private key that can be used for authentication on the resolver
func generateSignature(info *Info, privKey bmcrypto.PrivKey) string {
	// Generate token
	hash := sha256.Sum256([]byte(info.Hash + info.Routing))
	signature, err := bmcrypto.Sign(privKey, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}
