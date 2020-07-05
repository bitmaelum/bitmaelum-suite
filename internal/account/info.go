package account

import pow "github.com/bitmaelum/bitmaelum-server/pkg/proofofwork"

// Info represents client account information
type Info struct {
	Address string `json:"address"` // The address of the account

	// TODO: here be additional information about the user.
	Name         string `json:"name"`         // Full name of the user
	Organisation string `json:"organisation"` // Org of the user (if any)

	// Communication and encryption information
	PrivKey string          `json:"privKey"` // PEM encoded private key
	PubKey  string          `json:"pubKey"`  // PEM encoded public key
	Pow     pow.ProofOfWork `json:"pow"`     // Proof of work
	Server  string          `json:"server"`  // Mail server hosting this account
}
