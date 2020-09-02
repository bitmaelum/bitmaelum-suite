package resolve

import (
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

var errKeyNotFound = errors.New("hash not found")

// Repository is the interface to manage address resolving
type Repository interface {
	Resolve(addr address.HashAddress) (*Info, error)
	Upload(info *Info, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error
	Delete(info *Info, privKey bmcrypto.PrivKey) error
}
