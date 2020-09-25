package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// ChainRepository holds a list of multiple repositories which can all be tried to resolve addresses and keys
type ChainRepository struct {
	repos []Repository
}

// NewChainRepository Return a new chain repository
func NewChainRepository() *ChainRepository {
	return &ChainRepository{}
}

// Add a new repository to the chain
func (r *ChainRepository) Add(repo Repository) error {
	r.repos = append(r.repos, repo)

	return nil
}

// ResolveAddress an address through the chained repos
func (r *ChainRepository) ResolveAddress(addr address.HashAddress) (*AddressInfo, error) {
	for idx := range r.repos {
		info, err := r.repos[idx].ResolveAddress(addr)
		if err == nil {
			return info, nil
		}
	}

	return nil, errKeyNotFound
}

// ResolveRouting resolves routing
func (r *ChainRepository) ResolveRouting(routingID string) (*RoutingInfo, error) {
	for idx := range r.repos {
		info, err := r.repos[idx].ResolveRouting(routingID)
		if err == nil {
			return info, nil
		}
	}

	return nil, errKeyNotFound
}

// ResolveOrganisation resolves organisation
func (r *ChainRepository) ResolveOrganisation(orgHash string) (*OrganisationInfo, error) {
	for idx := range r.repos {
		info, err := r.repos[idx].ResolveOrganisation(orgHash)
		if err == nil {
			return info, nil
		}
	}

	return nil, errKeyNotFound
}

// UploadAddress public key through the chained repos
func (r *ChainRepository) UploadAddress(info *AddressInfo, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	for idx := range r.repos {
		err := r.repos[idx].UploadAddress(info, privKey, pow)
		if err != nil {
			return err
		}
	}

	return nil
}

// UploadRouting uploads routing information
func (r *ChainRepository) UploadRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	for idx := range r.repos {
		err := r.repos[idx].UploadRouting(info, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// UploadOrganisation uploads organisation information
func (r *ChainRepository) UploadOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	for idx := range r.repos {
		err := r.repos[idx].UploadOrganisation(info, privKey, pow)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAddress from repos
func (r *ChainRepository) DeleteAddress(info *AddressInfo, privKey bmcrypto.PrivKey) error {
	for idx := range r.repos {
		err := r.repos[idx].DeleteAddress(info, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteRouting from repos
func (r *ChainRepository) DeleteRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	for idx := range r.repos {
		err := r.repos[idx].DeleteRouting(info, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteOrganisation from repos
func (r *ChainRepository) DeleteOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey) error {
	for idx := range r.repos {
		err := r.repos[idx].DeleteOrganisation(info, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}
