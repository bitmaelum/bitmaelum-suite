package invite

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"time"
)

// Repository is the generic repository for dealing with invitations
type Repository interface {
	CreateInvite(addr address.HashAddress, expiry time.Duration) (string, error)
	GetInvite(addr address.HashAddress) (string, error)
	RemoveInvite(addr address.HashAddress) error
}
