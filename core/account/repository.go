package account

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
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
	GetBox(addr address.HashAddress, box string) (*message.MailBoxInfo, error)
	FindBox(addr address.HashAddress, query string) ([]message.MailBoxInfo, error)

	SendToBox(addr address.HashAddress, box, msgID string) error
	MoveToBox(addr address.HashAddress, srcBox, dstBox, msgID string) error

	// Message boxes
	FetchListFromBox(addr address.HashAddress, box string, offset, limit int) ([]message.List, error)

	// Flags
	GetFlags(addr address.HashAddress, box string, id string) ([]string, error)
	SetFlag(addr address.HashAddress, box string, id string, flag string) error
	UnsetFlag(addr address.HashAddress, box string, id string, flag string) error
}
