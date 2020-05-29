package resolve

import (
    "errors"
)

type ChainRepository struct {
    repos []Repository
}

func NewChainRepository() *ChainRepository {
    return &ChainRepository{}
}

func (r *ChainRepository) Add(repo Repository) error {
    r.repos = append(r.repos, repo)

    return nil
}

func (r *ChainRepository) Retrieve(hash string) (*ResolveInfo, error) {
    for idx := range(r.repos) {
        info, err := r.repos[idx].Retrieve(hash)
        if err == nil {
            return info, nil
        }
    }

    return nil, errors.New("Resolve info not found")
}

func (r *ChainRepository) Upload(hash, pubKey, address, signature string) error {
    for idx := range(r.repos) {
        err := r.repos[idx].Upload(hash, pubKey, address, signature)
        if err != nil {
            return err
        }
    }

    return nil
}
