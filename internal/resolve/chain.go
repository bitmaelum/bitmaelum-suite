package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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
func (r *ChainRepository) Upload(addr address.HashAddress, pubKey encrypt.PubKey, address, signature string) error {
	for idx := range r.repos {
		err := r.repos[idx].Upload(addr, pubKey, address, signature)
		if err != nil {
			return err
		}
	}

	return nil
}
