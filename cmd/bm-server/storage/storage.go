package storage

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/google/uuid"
	"io"
	"time"
)

type ProofOfWork struct {
	Challenge   string      `json:"challenge"`
	Bits        int         `json:"bits"`
	Nonce       uint64      `json:"nonce"`
	Expires     time.Time   `json:"expires"`
	MsgID       string      `json:"msg_id"`
	Valid       bool        `json:"valid"`
}

func (pow *ProofOfWork) MarshalBinary() (data []byte, err error) {
	return json.Marshal(pow)
}
func (pow *ProofOfWork) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, pow)
}


type Storable interface {
	// Retrieves the given challenge and returns it's proof of work info
	Retrieve(challenge string) (*ProofOfWork, error)
	// Stores the given proof of work
	Store(pow *ProofOfWork) error
	// Removes the given challenge
	Remove(challenge string) error
}

// NewProofOfWork generates a new proof of work
func NewProofOfWork() (*ProofOfWork, error) {
	// Generate a challenge the requesting server needs to validate
	challengeBuf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, challengeBuf)
	if err != nil {
		return nil, err
	}
	challenge := base64.StdEncoding.EncodeToString(challengeBuf)

	// Generate msgID we send back to the requestor
	tmp, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	msgID := tmp.String()

	// Store proof-of-work challenge into Redis
	pow := &ProofOfWork{
		Challenge: challenge,
		Bits:      config.Server.Accounts.ProofOfWork,
		Expires:   time.Now().Add(30 * time.Minute),
		Valid:     false,
		MsgID:     msgID,
	}

	return pow, nil
}
