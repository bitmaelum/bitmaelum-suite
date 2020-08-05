package account

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"time"
)

// Message is a simple message structure that we return as a list
type Message struct {
	Header  message.Header `json:"h"`
	Catalog []byte `json:"c"`
}

// MessageList is a list of messages we return
type MessageList struct {
	Total    int       `json:"total"`
	Returned int       `json:"returned"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Messages []Message `json:"messages"`
}

// Repository is an interface to manage accounts on a REMOTE machine (ie: server, not client side)
// @TODO: Too big.. Split into smaller interfaces
type Repository interface {
	// Account management
	Create(addr address.HashAddress) error
	Exists(addr address.HashAddress) bool

	// Public key
	StorePubKey(addr address.HashAddress, key string) error
	FetchPubKeys(addr address.HashAddress) ([]string, error)

	// Box related functions
	CreateBox(addr address.HashAddress, parentBox int) error
	ExistsBox(addr address.HashAddress, box int) bool
	DeleteBox(addr address.HashAddress, box int) error
	// GetBox(addr address.HashAddress, box int) (*message.MailBoxInfo, error)
	GetAllBoxes(addr address.HashAddress) ([]BoxInfo, error)
	//	FindBox(addr address.HashAddress, query string) ([]message.MailBoxInfo, error)

	SendToBox(addr address.HashAddress, box int, msgID string) error
	MoveToBox(addr address.HashAddress, srcBox, dstBox int, msgID string) error

	// Message boxes
	FetchListFromBox(addr address.HashAddress, box int, since time.Time, offset, limit int) (*MessageList, error)

	// Flags
	// @TODO flag repository? Are we going to use flags on messages this way?
	GetFlags(addr address.HashAddress, box int, id string) ([]string, error)
	SetFlag(addr address.HashAddress, box int, id string, flag string) error
	UnsetFlag(addr address.HashAddress, box int, id string, flag string) error
}
