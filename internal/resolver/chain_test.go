package resolver

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
