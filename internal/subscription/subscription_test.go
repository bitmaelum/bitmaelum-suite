package subscription

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscription(t *testing.T) {
	from, _ := address.NewHash("foo!")
	to, _ := address.NewHash("bar!")
	sub := New(*from, *to, "foobar")

	assert.Equal(t, *from, sub.From)
	assert.Equal(t, *to, sub.To)
	assert.Equal(t, "foobar", sub.SubscriptionID)
}
