package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplates(t *testing.T) {
	var buf = bytes.Buffer{}
	err := GenerateServerConfig(&buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())

	assert.Empty(t, Server.Logging.Level)
	err = Server.LoadConfig(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "trace", Server.Logging.Level)

	err = GenerateClientConfig(&buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())

	assert.Empty(t, Client.Resolver.Remote.URL)
	err = Client.LoadConfig(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "https://resolver.bitmaelum.com", Client.Resolver.Remote.URL)
}
