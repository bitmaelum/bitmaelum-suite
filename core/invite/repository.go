package invite

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"time"
)

// Repository is the generic repository for dealing with invitations
type Repository interface {
	CreateInvite(addr core.HashAddress, expiry time.Duration) (string, error)
	GetInvite(addr core.HashAddress) (string, error)
	RemoveInvite(addr core.HashAddress) error
}
