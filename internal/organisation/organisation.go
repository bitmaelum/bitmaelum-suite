package organisation

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// Organisation is a structure that defines an organiation
type Organisation struct {
	Addr       address.OrganisationHash `json:"address"`
	Name       string                   `json:"name"`
	PublicKey  bmcrypto.PubKey          `json:"public_key"`
	Validation []ValidationType         `json:"validations"`
}
