package pkg

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// Info represents client account information
type Info struct {
	Default bool   `json:"default"` // Is this the default account
	Address string `json:"address"` // The address of the account

	Name     string            `json:"name"`     // Full name of the user
	Settings map[string]string `json:"settings"` // Additional settings that can be user-defined

	// Communication and encryption information
	PrivKey encrypt.PrivKey `json:"privKey"` // PEM encoded private key
	PubKey  encrypt.PubKey  `json:"pubKey"`  // PEM encoded public key
	Pow     pow.ProofOfWork `json:"pow"`     // Proof of work
	Server  string          `json:"server"`  // Mail server hosting this account
}
