package resolve

import (
	account2 "github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

type localRepo struct {
	accountRepo account2.AddressRepository
}

// NewLocalRepository intializes a local repository
func NewLocalRepository(ar account2.AddressRepository) Repository {
	return &localRepo{
		accountRepo: ar,
	}
}

// Resolve resolves a local address
func (r *localRepo) Resolve(addr address.HashAddress) (*Info, error) {
	return nil, errKeyNotFound
}

func (r *localRepo) Upload(addr address.HashAddress, pubKey, address, signature string) error {
	return nil
}
