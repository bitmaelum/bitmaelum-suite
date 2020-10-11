package ticket

import (
	"os"
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func Test_BoltStorage(t *testing.T) {
	path := os.TempDir()
	b := NewBoltRepository(&path)
	assert.NotNil(t, b)

	from := hash.New("foo!")
	to := hash.New("bar!")
	tckt := NewUnvalidated(from, to, "foobar")

	err := b.Store(tckt)
	assert.Nil(t, err)

	tckt2, err := b.Fetch(tckt.ID)
	assert.Nil(t, err)
	assert.Equal(t, tckt2, tckt)

	b.Remove(tckt.ID)

	_, err = b.Fetch(tckt.ID)
	assert.NotNil(t, err)
}
