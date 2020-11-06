// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proof_of_work

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/google/uuid"
)

// ProofOfWork is the structure that keeps information about proof-of-work done for incoming messages. It connects
// the proof-of-work with a message ID which can be used for uploading.
type ProofOfWork struct {
	Challenge string    `json:"challenge"`
	Bits      int       `json:"bits"`
	Proof     uint64    `json:"proof"`
	Expires   time.Time `json:"expires"`
	MsgID     string    `json:"msg_id"`
	Valid     bool      `json:"valid"`
}

// Storable interface is the main interface to store and retrieve proof-of-work
type Storable interface {
	// Retrieve retrieves the given challenge from the storage and returns its proof-of-work info
	Retrieve(challenge string) (*ProofOfWork, error)
	// Store stores the given proof of work in the storage
	Store(pow *ProofOfWork) error
	// Remove removes the given challenge from the storage
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
