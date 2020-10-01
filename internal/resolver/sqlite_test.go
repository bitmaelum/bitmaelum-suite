package resolver

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqLiteAddress(t *testing.T) {
	r, err := NewSqliteRepository(":memory:")
	assert.NoError(t, err)

	testRepoAddress(t, r)
}

func TestSqLiteRouting(t *testing.T) {
	r, err := NewSqliteRepository(":memory:")
	assert.NoError(t, err)

	testRepoRouting(t, r)
}

func TestSqLiteOrganisation(t *testing.T) {
	r, err := NewSqliteRepository(":memory:")
	assert.NoError(t, err)

	testRepoOrganisation(t, r)
}
