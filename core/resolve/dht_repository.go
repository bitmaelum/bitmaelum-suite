package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/sirupsen/logrus"
)

type dhtRepo struct {
}

// NewDHTRepository generates a new DHT repository
func NewDHTRepository() Repository {
	return &dhtRepo{}
}

func (r *dhtRepo) Resolve(addr core.HashAddress) (*Info, error) {
	logrus.Trace("DHT is not implemented yet")

	return nil, errors.New("key not found in DHT")
}

func (r *dhtRepo) Upload(addr core.HashAddress, pubKey, address, signature string) error {
	return nil
}
