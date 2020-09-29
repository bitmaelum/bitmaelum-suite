package resolver

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
)

// Service represents a resolver service tied to a specific repository
type Service struct {
	repo         Repository
	routingCache *lru.Cache
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
	Hash        string                        `json:"hash"`        // Hash of the organisation
	PublicKey   bmcrypto.PubKey               `json:"public_key"`  // PublicKey of the organisation
	Pow         string                        `json:"pow"`         // Proof of work
	Validations []organisation.ValidationType `json:"validations"` // Validations for this organisation
}

// KeyRetrievalService initialises a key retrieval service.
func KeyRetrievalService(repo Repository) *Service {
	c, err := lru.New(64)
	if err != nil {
		c = nil
	}

	return &Service{
		repo:         repo,
		routingCache: c,
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
	// Fetch from cache if available
	if s.routingCache != nil {
		info, ok := s.routingCache.Get(routingID)
		if ok {
			logrus.Debugf("Resolving cached %s", routingID)
			return info.(*RoutingInfo), nil
		}
	}

	logrus.Debugf("Resolving %s", routingID)
	info, err := s.repo.ResolveRouting(routingID)
	if err != nil {
		logrus.Debugf("Error while resolving route %s: %s", routingID, err)
		return nil, err
	}

	// Store in cache if available
	if s.routingCache != nil {
		_ = s.routingCache.Add(routingID, info)
	}
	return info, nil
}

// ResolveOrganisation resolves a route.
func (s *Service) ResolveOrganisation(orgHash address.HashOrganisation) (*OrganisationInfo, error) {
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

// UploadOrganisationInfo uploads resolve information to one (or more) resolvers
func (s *Service) UploadOrganisationInfo(info internal.OrganisationInfo) error {
	a, err := address.NewOrgHash(info.Addr)
	if err != nil {
		return err
	}

	return s.repo.UploadOrganisation(&OrganisationInfo{
		Hash:        a.String(),
		PublicKey:   info.PubKey,
		Pow:         info.Pow.String(),
		Validations: info.Validations,
	}, info.PrivKey, info.Pow)
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
