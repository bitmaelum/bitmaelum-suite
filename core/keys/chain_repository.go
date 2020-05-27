package keys

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

func (r *ChainRepository) Retrieve(hash string) (string, error) {
    for idx := range(r.repos) {
        pubKey, err := r.repos[idx].Retrieve(hash)
        if err == nil {
            return pubKey, nil
        }
    }

    return "", errors.New("Key found found")
}
