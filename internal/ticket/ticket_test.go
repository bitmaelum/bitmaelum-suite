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
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestTicket(t *testing.T) {
	config.Server.Accounts.ProofOfWork = 4

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := NewUnvalidated(from, to, "foobar")

	assert.Equal(t, from, tckt.From)
	assert.Equal(t, to, tckt.To)
	assert.Equal(t, "foobar", tckt.SubscriptionID)
	assert.NotEmpty(t, tckt.ID)
	assert.False(t, tckt.Valid)
	assert.Equal(t, 4, tckt.Proof.Bits)
	assert.NotEmpty(t, tckt.Proof.Data)
}

func TestValidTicket(t *testing.T) {
	config.Server.Accounts.ProofOfWork = 4

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := NewValidated(from, to, "foobar")

	assert.Equal(t, from, tckt.From)
	assert.Equal(t, to, tckt.To)
	assert.Equal(t, "foobar", tckt.SubscriptionID)
	assert.NotEmpty(t, tckt.ID)
	assert.True(t, tckt.Valid)
}

func TestNewSimpleTicket(t *testing.T) {
	config.Server.Accounts.ProofOfWork = 4

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := NewValidated(from, to, "foobar")

	tckt2 := NewSimpleTicket(tckt)
	assert.Equal(t, tckt.ID, tckt2.ID)
	assert.Equal(t, tckt.Proof, tckt2.Proof)
	assert.Equal(t, tckt.Valid, tckt2.Valid)
}

func TestCreateTicketId(t *testing.T) {
	assert.Equal(t, "ticket-foo", createTicketKey("foo"))
}
