package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueIncomingMessage(t *testing.T) {
	IncomingChannel = make(chan string, 1)

	QueueIncomingMessage("12345")

	id := <- IncomingChannel
	assert.Equal(t, "12345", id)
}
