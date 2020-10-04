package subscription

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestSubscription(t *testing.T) {
	from := hash.New("foo!")
	to := hash.New("bar!")
	sub := New(from, to, "foobar")

	assert.Equal(t, from, sub.From)
	assert.Equal(t, to, sub.To)
	assert.Equal(t, "foobar", sub.SubscriptionID)
}

func TestCreateKey(t *testing.T) {
	from := hash.New("foo!")
	to := hash.New("bar!")
	sub := New(from, to, "foobar")

	assert.Equal(t, "sub-a6ca63d14d1c6c31ab71f60e7cd453aeac441e78372cddaa19667c05e45761e8", createKey(&sub))
}
