package resolve

import (
    "errors"
    "github.com/jaytaph/mailv2/core/account"
    "github.com/sirupsen/logrus"
)

type localRepo struct {
    as *account.Service
}

func NewLocalRepository(s *account.Service) Repository {
    return &localRepo{
        as: s,
    }
}

func (r *localRepo) Retrieve(hash string) (*ResolveInfo, error) {
    logrus.Trace("local repository cache is not available.")

    return nil, errors.New("key not found in local cache")
}

func (r *localRepo) Upload(hash, pubKey, address, signature string) error {
    return nil
}
