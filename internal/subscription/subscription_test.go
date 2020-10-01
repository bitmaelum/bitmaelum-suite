package subscription

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	from, _ := address.NewHash("foo!")
	to, _ := address.NewHash("bar!")
	sub := New(*from, *to, "foobar")

	assert.Equal(t, *from, sub.From)
	assert.Equal(t, *to, sub.To)
	assert.Equal(t, "foobar", sub.SubscriptionID)
}

func TestCreateKey(t *testing.T) {
	from, _ := address.NewHash("foo!")
	to, _ := address.NewHash("bar!")
	sub := New(*from, *to, "foobar")

	assert.Equal(t, "sub-40ab5cfb7bee2e2f2eb3e7b05a83ecba03a82a7920a182e614c0bce67602bfea", createKey(&sub))
}
