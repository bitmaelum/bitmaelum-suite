package ticket

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/google/uuid"
)

// TicketHeader is the HTTP Header that contains our ticket ID
const TicketHeader = "x-bitmaelum-ticket"

// SimpleTicket is a structure that holds id, proof of work and if it's valid or not. Used for output because normally
// we use Ticket instead.
type SimpleTicket struct {
	ID    string           `json:"ticket_id"`     // ticket ID. Will be used as the message ID when uploading
	Pow   *pow.ProofOfWork `json:"proof_of_work"` // proof of work that must be completed
	Valid bool             `json:"is_valid"`      // true if the ticket is valid
}

// Ticket is a structure that defines if a client or server is allowed to upload a message, or if additional work has to be done first
type Ticket struct {
	ID             string              `json:"ticket_id"`       // ticket ID. Will be used as the message ID when uploading
	Pow            *pow.ProofOfWork    `json:"proof_of_work"`   // proof of work that must be completed
	Valid          bool                `json:"is_valid"`        // true if the ticket is valid
	From           address.HashAddress `json:"from_addr"`       // From address for this ticket
	To             address.HashAddress `json:"to_addr"`         // To address for this ticket
	SubscriptionID string              `json:"subscription_id"` // mailing list subscription ID (if any)
}

// MarshalBinary converts a ticket to binary format so it can be stored in Redis
func (t *Ticket) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UnmarshalBinary converts binary to a ticket so it can be fetched from Redis
func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// New creates a new unvalidated ticket with proof of work
func New(from, to address.HashAddress, subscriptionID string) *Ticket {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil
	}

	proof := pow.New(config.Server.Accounts.ProofOfWork, pow.GenerateWork(), 0)
	t := &Ticket{
		ID:             id.String(),
		Pow:            proof,
		Valid:          false,
		From:           from,
		To:             to,
		SubscriptionID: subscriptionID,
	}

	return t
}

// NewValid returns a new ticket that is validated (without proof-of-work)
func NewValid(from, to address.HashAddress, subscriptionID string) *Ticket {
	t := New(from, to, subscriptionID)
	t.Valid = true

	return t
}

// NewSimpleTicket converts a ticket into a simple ticket. Used for outputting when we don't need routing info
func NewSimpleTicket(t *Ticket) SimpleTicket {
	return SimpleTicket{
		ID:    t.ID,
		Pow:   t.Pow,
		Valid: t.Valid,
	}
}

// Repository is a ticket repository to fetch and store tickets
type Repository interface {
	Fetch(ticketID string) (*Ticket, error)
	Store(ticket *Ticket) error
	Remove(ticketID string)
}
