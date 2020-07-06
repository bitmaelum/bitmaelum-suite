package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/core/account"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

type localRepo struct {
	as *account.Service
}

// NewLocalRepository intializes a local repository
func NewLocalRepository(s *account.Service) Repository {
	return &localRepo{
		as: s,
	}
}

// Resolve resolves a local address
func (r *localRepo) Resolve(addr address.HashAddress) (*Info, error) {
	logrus.Trace("local repository cache is not available.")

	return nil, errors.New("key not found in local cache")
}

func (r *localRepo) Upload(addr address.HashAddress, pubKey, address, signature string) error {
	return nil
}
