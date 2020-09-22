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

// AddressInfo is a structure returned by the external resolver system
type AddressInfo struct {
	Hash        string          `json:"hash"`       // Hash of the email address
	PublicKey   bmcrypto.PubKey `json:"public_key"` // PublicKey of the user
	RoutingID   string          `json:"routing"`    // Routing ID
	Pow         string          `json:"pow"`        // Proof of work
	RoutingInfo RoutingInfo     `json:"_"`          // Don't store
}

// RoutingInfo is a structure returned by the external resolver system
type RoutingInfo struct {
	Hash      string          `json:"hash"`       // Hash / routingID
	PublicKey bmcrypto.PubKey `json:"public_key"` // PublicKey of the user
	Routing   string          `json:"routing"`    // Server where this email address resides
}

// OrganisationInfo is a structure returned by the external resolver system
type OrganisationInfo struct {
	Hash      string          `json:"hash"`       // Hash of the organisation
	PublicKey bmcrypto.PubKey `json:"public_key"` // PublicKey of the organisation
	Pow       string          `json:"pow"`        // Proof of work
}

// KeyRetrievalService initialises a key retrieval service.
func KeyRetrievalService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ResolveAddress resolves an address.
func (s *Service) ResolveAddress(addr address.HashAddress) (*AddressInfo, error) {
	logrus.Debugf("Resolving address %s", addr.String())
	info, err := s.repo.ResolveAddress(addr)
	if err != nil {
		logrus.Debugf("Error while resolving address %s: %s", addr.String(), err)
		return nil, err
	}

	// Resolve routing info, we often need it
	ri, err := s.repo.ResolveRouting(info.RoutingID)
	if err != nil {
		logrus.Debugf("Error while resolving routing %s for address %s: %s", info.RoutingID, addr.String(), err)
		return nil, err
	}
	info.RoutingInfo = *ri

	return info, err
}

// ResolveRouting resolves a route.
func (s *Service) ResolveRouting(routingID string) (*RoutingInfo, error) {
	logrus.Debugf("Resolving %s", routingID)
	info, err := s.repo.ResolveRouting(routingID)
	if err != nil {
		logrus.Debugf("Error while resolving route %s: %s", routingID, err)
	}

	return info, err
}

// ResolveOrganisation resolves a route.
func (s *Service) ResolveOrganisation(orgHash string) (*OrganisationInfo, error) {
	logrus.Debugf("Resolving %s", orgHash)
	info, err := s.repo.ResolveOrganisation(orgHash)
	if err != nil {
		logrus.Debugf("Error while resolving organisation %s: %s", orgHash, err)
	}

	return info, err
}

// UploadAddressInfo uploads resolve information to one (or more) resolvers
func (s *Service) UploadAddressInfo(info internal.AccountInfo) error {
	hashAddr, err := address.NewHash(info.Address)
	if err != nil {
		return err
	}

	return s.repo.UploadAddress(&AddressInfo{
		Hash:      hashAddr.String(),
		PublicKey: info.PubKey,
		RoutingID: info.RoutingID,
		Pow:       info.Pow.String(),
	}, info.PrivKey, info.Pow)
}

// UploadRoutingInfo uploads resolve information to one (or more) resolvers
func (s *Service) UploadRoutingInfo(info internal.RoutingInfo) error {

	return s.repo.UploadRouting(&RoutingInfo{
		Hash:      info.RoutingID,
		PublicKey: info.PubKey,
		Routing:   info.Route,
	}, info.PrivKey)
}

// generateAddressSignature generates a signature with the accounts private key that can be used for authentication on the resolver
func generateAddressSignature(info *AddressInfo, privKey bmcrypto.PrivKey) string {
	// Generate token
	hash := sha256.Sum256([]byte(info.Hash + info.RoutingID))
	signature, err := bmcrypto.Sign(privKey, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

// generateRoutingSignature generates a signature with the accounts private key that can be used for authentication on the resolver
func generateRoutingSignature(info *RoutingInfo, privKey bmcrypto.PrivKey) string {
	// Generate token
	hash := sha256.Sum256([]byte(info.Hash))
	signature, err := bmcrypto.Sign(privKey, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

// generateOrganisationSignature generates a signature with the accounts private key that can be used for authentication on the resolver
func generateOrganisationSignature(info *OrganisationInfo, privKey bmcrypto.PrivKey) string {
	// Generate token
	hash := sha256.Sum256([]byte(info.Hash))
	signature, err := bmcrypto.Sign(privKey, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}
