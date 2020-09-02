package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
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

func (r *dhtRepo) Upload(info *Info, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	return nil
}

func (r *dhtRepo) Delete(info *Info, privKey bmcrypto.PrivKey) error {
	return nil
}
