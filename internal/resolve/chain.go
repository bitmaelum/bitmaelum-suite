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

// Resolve an address through the chained repos
func (r *ChainRepository) Resolve(addr address.HashAddress) (*Info, error) {
	for idx := range r.repos {
		info, err := r.repos[idx].Resolve(addr)
		if err == nil {
			return info, nil
		}
	}

	return nil, errKeyNotFound
}

// Upload public key through the chained repos
func (r *ChainRepository) Upload(info *Info, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	for idx := range r.repos {
		err := r.repos[idx].Upload(info, privKey, pow)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete from repos
func (r *ChainRepository) Delete(info *Info, privKey bmcrypto.PrivKey) error {
	for idx := range r.repos {
		err := r.repos[idx].Delete(info, privKey)
		if err != nil {
			return err
		}
	}

	return nil
}
