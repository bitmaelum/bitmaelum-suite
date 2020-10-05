package organisation

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// Organisation is a structure that defines an organiation
type Organisation struct {
	Addr       hash.Hash        `json:"address"`
	Name       string           `json:"name"`
	PublicKey  bmcrypto.PubKey  `json:"public_key"`
	Validation []ValidationType `json:"validations"`
}
