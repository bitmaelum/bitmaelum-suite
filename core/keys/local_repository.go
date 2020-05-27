package keys

import (
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/api"
)

type localRepo struct {
    as *account.Service
}

func NewLocalRepository(s *account.Service) Repository {
    return &localRepo{
        as: s,
    }
}

func (r *localRepo) Retrieve(hash string) (string, error) {
    client, err := api.NewClient("127.0.0.1", 2424)
    if err != nil {
        return "", err
    }

    return client.GetPublicKey(hash)
}
