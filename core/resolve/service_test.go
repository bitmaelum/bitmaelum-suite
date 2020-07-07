package resolve

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsLocal(t *testing.T) {
	info := &Info{
		Hash:      "12345",
		PublicKey: "FOOBAR",
		Server:    "my.host:2424",
	}

	config.Server.Server.Name = "my.host:2424"
	assert.True(t, info.IsLocal())

	config.Server.Server.Name = "my.host"
	assert.True(t, info.IsLocal())

	config.Server.Server.Name = "my.host:2425"
	assert.False(t, info.IsLocal())

	config.Server.Server.Name = "another.host:2425"
	assert.False(t, info.IsLocal())

	config.Server.Server.Name = "https://my.host:2424"
	assert.True(t, info.IsLocal())

	config.Server.Server.Name = "http://my.host:2424"
	assert.True(t, info.IsLocal())
}
