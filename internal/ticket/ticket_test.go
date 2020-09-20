package ticket

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTicket(t *testing.T) {
	config.Server.Accounts.ProofOfWork = 22

	from, _ := address.NewHash("foo!")
	to, _ := address.NewHash("bar!")
	tckt := NewUnvalidated(*from, *to, "foobar")

	assert.Equal(t, *from, tckt.From)
	assert.Equal(t, *to, tckt.To)
	assert.Equal(t, "foobar", tckt.SubscriptionID)
	assert.NotEmpty(t, tckt.ID)
	assert.False(t, tckt.Valid)
	assert.Equal(t, 22, tckt.Proof.Bits)
	assert.NotEmpty(t, tckt.Proof.Data)
}

func TestValidTicket(t *testing.T) {
	config.Server.Accounts.ProofOfWork = 22

	from, _ := address.NewHash("foo!")
	to, _ := address.NewHash("bar!")
	tckt := NewValidated(*from, *to, "foobar")

	assert.Equal(t, *from, tckt.From)
	assert.Equal(t, *to, tckt.To)
	assert.Equal(t, "foobar", tckt.SubscriptionID)
	assert.NotEmpty(t, tckt.ID)
	assert.True(t, tckt.Valid)
}
