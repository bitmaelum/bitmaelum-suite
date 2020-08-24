package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

type dhtRepo struct {
}

// NewDHTRepository generates a new DHT repository
func NewDHTRepository() Repository {
	return &dhtRepo{}
}

func (r *dhtRepo) Resolve(addr address.HashAddress) (*Info, error) {
	return nil, errKeyNotFound
}

func (r *dhtRepo) Upload(addr address.HashAddress, pubKey, address, signature string) error {
	return nil
}
