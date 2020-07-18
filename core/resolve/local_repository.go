package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/core/account"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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
	return nil, errKeyNotFound
}

func (r *localRepo) Upload(addr address.HashAddress, pubKey, address, signature string) error {
	return nil
}
