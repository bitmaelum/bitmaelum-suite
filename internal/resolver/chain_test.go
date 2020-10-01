package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainRepository(t *testing.T) {
	r := NewChainRepository()

	m, err := NewMockRepository()
	assert.NoError(t, err)
	err = r.Add(m)
	assert.NoError(t, err)

	testRepoAddress(t, r)
	testRepoRouting(t, r)
	testRepoOrganisation(t, r)
}
