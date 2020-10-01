package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockRepository(t *testing.T) {
	r, err := NewMockRepository()
	assert.NoError(t, err)

	testRepoAddress(t, r)
	testRepoRouting(t, r)
	testRepoOrganisation(t, r)
}
