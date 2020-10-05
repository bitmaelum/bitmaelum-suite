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
