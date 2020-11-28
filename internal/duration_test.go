package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationCorrect(t *testing.T) {
	dataProvider := []struct {
			s  string
			ms int64
		}{
		{s: "1d", ms: 24 * 3600 * 1000},
		{s: "10d", ms: 240 * 3600 * 1000},
		{s: "10d1h", ms: 241 * 3600 * 1000},
		{s: "10d1h1m", ms: (10 * 24 * 3600 + 1 * 3600 + 60) * 1000},
		{s: "1d1y", ms: 31622400000},
		{s: "1d1m", ms: 86460000},
		{s: "1m1d", ms: 86460000},
		{s: "5w141m", ms: 3032460000},
	}

	for _, entry := range dataProvider {
		d, err := ParseDuration(entry.s)
		assert.NoError(t, err)
		assert.Equal(t, entry.ms, d.Milliseconds())
	}
}

func TestDurationIncorrect(t *testing.T) {
	dataProvider := []string{
		"1h 1h",
		"",
		"1A",
		"-151w",
		"0w",
		"foobar",
		"5P1G",
		"1d1w4d",
	}

	for _, entry := range dataProvider {
		d, err := ParseDuration(entry)
		assert.Error(t, err)
		assert.Equal(t, d, time.Duration(0))
	}
}
