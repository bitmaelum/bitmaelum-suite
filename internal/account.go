package internal

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// AccountInfo represents client account information
type AccountInfo struct {
	Default bool   `json:"default"` // Is this the default account
	Address string `json:"address"` // The address of the account

	Name     string            `json:"name"`     // Full name of the user
	Settings map[string]string `json:"settings"` // Additional settings that can be user-defined

	// Communication and encryption information
	PrivKey bmcrypto.PrivKey        `json:"privKey"`       // PEM encoded private key
	PubKey  bmcrypto.PubKey         `json:"pubKey"`        // PEM encoded public key
	Pow     proofofwork.ProofOfWork `json:"pow,omitEmpty"` // Proof of work
	Routing string                  `json:"server"`        // Mail server hosting this account
}
