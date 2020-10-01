package processor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRetry(t *testing.T) {
	var d time.Duration

	d = getNextRetryDuration(0)
	assert.Equal(t, time.Duration(60000000000), d)

	d = getNextRetryDuration(1)
	assert.Equal(t, time.Duration(60000000000), d)

	d = getNextRetryDuration(8)
	assert.Equal(t, time.Duration(300000000000), d)

	d = getNextRetryDuration(16)
	assert.Equal(t, time.Duration(300000000000), d)

	d = getNextRetryDuration(17)
	assert.Equal(t, time.Duration(1800000000000), d)

	d = getNextRetryDuration(18)
	assert.Equal(t, time.Duration(1800000000000), d)

	d = getNextRetryDuration(31)
	assert.Equal(t, time.Duration(0), d)
}
