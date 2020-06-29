package account

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/messagebox"
)

// Repository is an interface to manage accounts on a REMOTE machine (ie: server, not client side)
type Repository interface {
	// Account management
	Create(addr core.HashAddress) error
	Exists(addr core.HashAddress) bool

	// Public key
	StorePubKey(addr core.HashAddress, key string) error
	FetchPubKeys(addr core.HashAddress) ([]string, error)

	// Box related functions
	CreateBox(addr core.HashAddress, box, name, description string, quota int) error
	ExistsBox(addr core.HashAddress, box string) bool
	DeleteBox(addr core.HashAddress, box string) error
	GetBox(addr core.HashAddress, box string) (*messagebox.MailBoxInfo, error)
	FindBox(addr core.HashAddress, query string) ([]messagebox.MailBoxInfo, error)

	// Message boxes
	FetchListFromBox(addr core.HashAddress, box string, offset, limit int) ([]messagebox.MessageList, error)

	// Flags
	GetFlags(addr core.HashAddress, box string, id string) ([]string, error)
	SetFlag(addr core.HashAddress, box string, id string, flag string) error
	UnsetFlag(addr core.HashAddress, box string, id string, flag string) error
}
