package organisation

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// Organisation is a structure that defines an organisation
type Organisation struct {
	Hash       hash.Hash        `json:"hash"`
	FullName   string           `json:"name"`
	PublicKey  bmcrypto.PubKey  `json:"public_key"`
	Validation []ValidationType `json:"validations"`
}
