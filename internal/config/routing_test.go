package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	r, err := Generate()
	assert.NoError(t, err)

	assert.NotEmpty(t, r.RoutingID)
	assert.NotEmpty(t, r.PrivateKey)
	assert.NotEmpty(t, r.PublicKey)
}

func TestReadSaveRouting(t *testing.T) {
	r, err := Generate()
	assert.NoError(t, err)

	fs = afero.NewMemMapFs()

	err = SaveRouting("/generated/routing.json", r)
	assert.NoError(t, err)

	err = ReadRouting("/generated/routing.json")
	assert.NoError(t, err)
	assert.Equal(t, r.RoutingID, Server.Routing.RoutingID)
}
