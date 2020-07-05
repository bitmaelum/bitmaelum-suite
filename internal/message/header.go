package message

import "github.com/bitmaelum/bitmaelum-server/pkg/address"
import pow "github.com/bitmaelum/bitmaelum-server/pkg/proofofwork"

// Header represents a message header
type Header struct {
	From struct {
		Addr        address.HashAddress `json:"address"`
		PublicKey   string              `json:"public_key"`
		ProofOfWork pow.ProofOfWork     `json:"proof_of_work"`
	} `json:"from"`
	To struct {
		Addr address.HashAddress `json:"address"`
	} `json:"to"`
	Catalog struct {
		Size         uint64     `json:"size"`
		Checksum     []Checksum `json:"checksum"`
		Crypto       string     `json:"crypto"`
		EncryptedKey []byte     `json:"encrypted_key"`
	} `json:"catalog"`
}

// Checksum holds a checksum which consists of the checksum hash value, and the given type of the checksum
type Checksum struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
}
