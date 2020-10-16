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

package ticket

import (
	"encoding/json"
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TicketHeader is the HTTP Header that contains our ticket ID
const TicketHeader = "x-bitmaelum-ticket"

// SimpleTicket is a structure that holds id, proof of work and if it's valid or not. Used for output because normally
// we use Ticket instead.
type SimpleTicket struct {
	ID    string           `json:"ticket_id"`     // ticket ID. Will be used as the message ID when uploading
	Proof *pow.ProofOfWork `json:"proof_of_work"` // proof of work that must be completed
	Valid bool             `json:"is_valid"`      // true if the ticket is valid
}

// Ticket is a structure that defines if a client or server is allowed to upload a message, or if additional work has to be done first
type Ticket struct {
	ID             string           `json:"ticket_id"`       // ticket ID. Will be used as the message ID when uploading
	Proof          *pow.ProofOfWork `json:"proof_of_work"`   // proof of work that must be completed
	Valid          bool             `json:"is_valid"`        // true if the ticket is valid
	From           hash.Hash        `json:"from_addr"`       // From address for this ticket
	To             hash.Hash        `json:"to_addr"`         // To address for this ticket
	SubscriptionID string           `json:"subscription_id"` // mailing list subscription ID (if any)
}

// MarshalBinary converts a ticket to binary format so it can be stored in Redis
func (t *Ticket) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

// UnmarshalBinary converts binary to a ticket so it can be fetched from Redis
func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// NewUnvalidated creates a new unvalidated ticket with proof of work
func NewUnvalidated(from, to hash.Hash, subscriptionID string) *Ticket {
	logrus.Trace("Generating new unvalidated ticket")

	// Generate Ticket ID
	ticketUUID, err := uuid.NewRandom()
	if err != nil {
		return nil
	}
	ticketID := ticketUUID.String()
	logrus.Trace("TicketID: ", ticketID)

	// Generate workdata for proof-of-work
	work, err := pow.GenerateWorkData()
	if err != nil {
		return nil
	}
	proof := pow.NewWithoutProof(config.Server.Accounts.ProofOfWork, work)

	// Return ticket
	return &Ticket{
		ID:             ticketID,
		Proof:          proof,
		Valid:          false,
		From:           from,
		To:             to,
		SubscriptionID: subscriptionID,
	}
}

// NewValidated returns a new ticket that is validated (without proof-of-work)
func NewValidated(from, to hash.Hash, subscriptionID string) *Ticket {
	tckt := NewUnvalidated(from, to, subscriptionID)
	tckt.Valid = true

	return tckt
}

// NewSimpleTicket converts a ticket into a simple ticket. Used for outputting when we don't need routing info
func NewSimpleTicket(t *Ticket) SimpleTicket {
	return SimpleTicket{
		ID:    t.ID,
		Proof: t.Proof,
		Valid: t.Valid,
	}
}

// Repository is a ticket repository to fetch and store tickets
type Repository interface {
	Fetch(ticketID string) (*Ticket, error)
	Store(ticket *Ticket) error
	Remove(ticketID string)
}

// createTicketKey creates a key based on the given ID. This is needed otherwise we might send any data as ticket-id
// to redis in order to extract other kind of data (and you don't want that).
func createTicketKey(id string) string {
	return fmt.Sprintf("ticket-%s", id)
}
