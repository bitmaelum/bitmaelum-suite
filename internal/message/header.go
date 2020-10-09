package message

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"

	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// ChecksumList is a list of key/value pairs of checksums. ie: ["sha1"] = "123456abcde"
type ChecksumList map[string]string

// Header represents a message header
type Header struct {
	From struct {
		Addr        hash.Hash        `json:"address"`
		PublicKey   *bmcrypto.PubKey `json:"public_key"`
		ProofOfWork *pow.ProofOfWork `json:"proof_of_work"`
	} `json:"from"`
	To struct {
		Addr hash.Hash `json:"address"`
	} `json:"to"`
	Catalog struct {
		Size          uint64       `json:"size"`
		Checksum      ChecksumList `json:"checksum"`
		Crypto        string       `json:"crypto"`
		TransactionID string       `json:"txID"`
		EncryptedKey  []byte       `json:"encrypted_key"`
	} `json:"catalog"`

	// Signature of the from, to and catalog structures, as signed by the private key of the server.
	ServerSignature string `json:"sender_signature,omitempty"`
}

// Checksum holds a checksum which consists of the checksum hash value, and the given type of the checksum
type Checksum struct {
	Hash  string `json:"hash"`
	Value string `json:"value"`
}
