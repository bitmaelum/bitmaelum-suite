package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {
	var err error

	err = MangementPermissions([]string{})
	assert.NoError(t, err)

	err = MangementPermissions([]string{"flush"})
	assert.NoError(t, err)

	err = MangementPermissions([]string{"flush", "invite"})
	assert.NoError(t, err)

	err = MangementPermissions([]string{"foo"})
	assert.Error(t, err)

	err = MangementPermissions([]string{"flush", "foo"})
	assert.Error(t, err)
}

func TestValidDuration(t *testing.T) {
	var d time.Duration
	var err error

	d, err = ValidDuration("")
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), d)

	_, err = ValidDuration("foobar")
	assert.Error(t, err)

	d, err = ValidDuration("1")
	assert.NoError(t, err)
	assert.Equal(t, 24*time.Hour, d)

	d, err = ValidDuration("1d")
	assert.NoError(t, err)
	assert.Equal(t, 24*time.Hour, d)

	d, err = ValidDuration("2h5m1s")
	assert.NoError(t, err)
	assert.Equal(t, 7501*time.Second, d)
}
