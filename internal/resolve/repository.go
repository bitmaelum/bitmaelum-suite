package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

var errKeyNotFound = errors.New("hash not found")

// Repository is a complete key resolver repository with the different parts
type Repository interface {
	AddressRepository
	RouteRepository
	OrganisationRepository
}

// AddressRepository is the interface to manage address resolving
type AddressRepository interface {
	ResolveAddress(addr address.HashAddress) (*AddressInfo, error)
	UploadAddress(info *AddressInfo, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error
	DeleteAddress(info *AddressInfo, privKey bmcrypto.PrivKey) error
}

// RouteRepository is the interface to manage route resolving
type RouteRepository interface {
	ResolveRouting(routingID string) (*RoutingInfo, error)
	UploadRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error
	DeleteRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error
}

// OrganisationRepository is the interface to manage organisation resolving
type OrganisationRepository interface {
	ResolveOrganisation(orgHash string) (*OrganisationInfo, error)
	UploadOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error
	DeleteOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey) error
}
