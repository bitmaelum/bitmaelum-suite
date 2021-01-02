// Copyright (c) 2021 BitMaelum Authors
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
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/work"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TicketHeader is the HTTP Header that contains our ticket ID when sending messages
const TicketHeader = "x-bitmaelum-ticket"

// WorkType is a structure that holds the type and work repository data
type WorkType struct {
	Type string          // Type of work
	Data work.Repository // Work stored on this ticket
}

// UnmarshalJSON We unmarshal from the ticketWorkType, as we do not know which repository implementation we need. However, we can
// marshal from the implementation. This is why there is no MarshalJSON present here.
func (twt *WorkType) UnmarshalJSON(data []byte) error {
	type tmpType struct {
		Type string
		Data string
	}

	tmp := &tmpType{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	twt.Type = tmp.Type

	switch tmp.Type {
	case "pow":
		twt.Data, err = work.NewPowFromString(tmp.Data)
	}

	return err
}

// Ticket is a structure that defines if a client or server is allowed to upload a message, or if additional work has to be done first
type Ticket struct {
	ID     string    // ticket ID
	Valid  bool      // true if the ticket is valid
	Expiry time.Time // Time when this ticket expires
	Work   *WorkType // Work that needs to be done for validating the ticket

	Sender         hash.Hash // From address for this ticket
	Recipient      hash.Hash // To address for this ticket
	SubscriptionID string    // mailing list subscription ID (if any)

	AuthKey string // Optional authkey attached to the ticket in case its send on behalf
}

// MarshalBinary converts a ticket to binary format so it can be stored in Redis
func (t *Ticket) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

// UnmarshalBinary converts binary to a ticket so it can be fetched from Redis
func (t *Ticket) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// Expired will return true when the ticket has expired
func (t *Ticket) Expired() bool {
	return t.Expiry.Before(internal.TimeNow())
}

// New creates a new unvalidated ticket without work
func New(senderHash, recipientHash hash.Hash, subscriptionID string) *Ticket {
	logrus.Trace("Generating new ticket")

	// Generate Ticket ID
	ticketUUID, err := uuid.NewRandom()
	if err != nil {
		return nil
	}
	ticketID := ticketUUID.String()
	logrus.Trace("TicketID: ", ticketID)

	// Return ticket
	return &Ticket{
		ID:             ticketID,
		Expiry:         internal.TimeNow().Add(1800 * time.Second),
		Valid:          false,
		Sender:         senderHash,
		Recipient:      recipientHash,
		SubscriptionID: subscriptionID,
	}
}

// NewFromBytes will generate a ticket from the json-encoded data
func NewFromBytes(data []byte) (*Ticket, error) {
	t := &Ticket{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
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
