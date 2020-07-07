package account

import (
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

// Repository is an interface to manage accounts on a REMOTE machine (ie: server, not client side)
type Repository interface {
	// Account management
	Create(addr address.HashAddress) error
	Exists(addr address.HashAddress) bool

	// Public key
	StorePubKey(addr address.HashAddress, key string) error
	FetchPubKeys(addr address.HashAddress) ([]string, error)

	// Box related functions
	CreateBox(addr address.HashAddress, box, name, description string, quota int) error
	ExistsBox(addr address.HashAddress, box string) bool
	DeleteBox(addr address.HashAddress, box string) error
	GetBox(addr address.HashAddress, box string) (*core.MailBoxInfo, error)
	FindBox(addr address.HashAddress, query string) ([]core.MailBoxInfo, error)

	// Message boxes
	FetchListFromBox(addr address.HashAddress, box string, offset, limit int) ([]core.MessageList, error)

	// Flags
	GetFlags(addr address.HashAddress, box string, id string) ([]string, error)
	SetFlag(addr address.HashAddress, box string, id string, flag string) error
	UnsetFlag(addr address.HashAddress, box string, id string, flag string) error
}
