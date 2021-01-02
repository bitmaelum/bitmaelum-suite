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
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/work"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestTicket(t *testing.T) {
	config.Server.Work.Pow.Bits = 4

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := New(from, to, "foobar")

	assert.Equal(t, from, tckt.Sender)
	assert.Equal(t, to, tckt.Recipient)
	assert.Equal(t, "foobar", tckt.SubscriptionID)
	assert.NotEmpty(t, tckt.ID)
	assert.False(t, tckt.Valid)

	workRepo, err := work.NewPow()
	assert.NoError(t, err)
	tckt.Work = &WorkType{
		Type: workRepo.GetName(),
		Data: workRepo,
	}

	out := tckt.Work.Data.GetWorkOutput()
	assert.Equal(t, 4, out["bits"])
	assert.NotEmpty(t, out["data"])
}

func TestCreateTicketId(t *testing.T) {
	assert.Equal(t, "ticket-foo", createTicketKey("foo"))
}

func TestTicketMarshalBinary(t *testing.T) {
	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := New(from, to, "foobar")

	b, err := tckt.MarshalBinary()
	assert.NoError(t, err)

	tckt = &Ticket{}
	err = tckt.UnmarshalBinary(b)
	assert.NoError(t, err)
	assert.Equal(t, "c0e0aaaea050bcf3be26c0c23d58fa890c0dfb79c8a23016b4a86cd28ca6ea71", tckt.Sender.String())
	assert.Equal(t, "e687b749f2cd93615923a2f705faace4033f35d57ccfca652cdc39616a94a3c2", tckt.Recipient.String())
	assert.Equal(t, "foobar", tckt.SubscriptionID)
}

func TestTicketMarshalling(t *testing.T) {
	config.Server.Work.Pow.Bits = 14

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := New(from, to, "")

	workRepo, err := work.NewPow()
	assert.NoError(t, err)
	tckt.Work = &WorkType{
		Type: workRepo.GetName(),
		Data: workRepo,
	}

	b, err := json.MarshalIndent(tckt, "", "  ")
	assert.NoError(t, err)

	tckt2 := &Ticket{}
	err = json.Unmarshal(b, &tckt2)
	assert.NoError(t, err)

	assert.Equal(t, tckt.Work.Data.GetWorkOutput(), tckt2.Work.Data.GetWorkOutput())
}
