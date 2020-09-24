package invite

import (
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Repository is the generic repository for dealing with invitations
type Repository interface {
	Create(addr address.HashAddress, expiry time.Duration) (string, error)
	Get(addr address.HashAddress) (string, error)
	Remove(addr address.HashAddress) error
}

func createInviteKey(addr address.HashAddress) string {
	return "invite." + addr.String()
}
