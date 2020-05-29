package resolve

import (
    "errors"
    "github.com/jaytaph/mailv2/core"
)

type ChainRepository struct {
    repos []Repository
}

// Return a new chain repository
func NewChainRepository() *ChainRepository {
    return &ChainRepository{}
}

// Add a new repository to the chain
func (r *ChainRepository) Add(repo Repository) error {
    r.repos = append(r.repos, repo)

    return nil
}

// Resolve an address through the chained repos
func (r *ChainRepository) Resolve(addr core.HashAddress) (*ResolveInfo, error) {
    for idx := range(r.repos) {
        info, err := r.repos[idx].Resolve(addr)
        if err == nil {
            return info, nil
        }
    }

    return nil, errors.New("Resolve info not found")
}

// Upload public key through the chained repos
func (r *ChainRepository) Upload(hash, pubKey, address, signature string) error {
    for idx := range(r.repos) {
        err := r.repos[idx].Upload(hash, pubKey, address, signature)
        if err != nil {
            return err
        }
    }

    return nil
}
