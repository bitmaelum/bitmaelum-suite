package account

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"time"
)

// Message is a simple message structure that we return as a list
type Message struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"h"`
	Catalog []byte         `json:"c"`
}

// MessageList is a list of messages we return
type MessageList struct {
	Total    int       `json:"total"`
	Returned int       `json:"returned"`
	Limit    int       `json:"limit"`
	Offset   int       `json:"offset"`
	Messages []Message `json:"messages"`
}

// BoxInfo returns information about the given message box
type BoxInfo struct {
	ID    int `json:"id"`
	Total int `json:"total"`
}

// Repository is the main repository that needs to be implemented. It's pretty big
type Repository interface {
	AddressRepository
	FlagRepository
	KeyRepository
	BoxRepository
	MessageRepository
}

// AddressRepository creates, checks or deletes complete accounts. Address is not the correct word for this.
type AddressRepository interface {
	Create(addr address.HashAddress, pubKey string) error
	Exists(addr address.HashAddress) bool
	Delete(addr address.HashAddress) error
}

// FlagRepository is for dealing with message flags
// @TODO flag repository? Are we going to use flags on messages this way?
type FlagRepository interface {
	// Flags
	GetFlags(addr address.HashAddress, box int, id string) ([]string, error)
	SetFlag(addr address.HashAddress, box int, id string, flag string) error
	UnsetFlag(addr address.HashAddress, box int, id string, flag string) error
}

// KeyRepository gets and sets public keys into an account
type KeyRepository interface {
	// Public key
	StoreKey(addr address.HashAddress, key string) error
	FetchKeys(addr address.HashAddress) ([]string, error)
	FetchDecodedKeys(addr address.HashAddress) ([]interface{}, error)
}

// BoxRepository deals with message boxes insides an account
type BoxRepository interface {
	CreateBox(addr address.HashAddress, parentBox int) error
	ExistsBox(addr address.HashAddress, box int) bool
	DeleteBox(addr address.HashAddress, box int) error
	GetBoxInfo(addr address.HashAddress, box int) (*BoxInfo, error)
	GetAllBoxes(addr address.HashAddress) ([]BoxInfo, error)
}

// MessageRepository deals with message within boxes
type MessageRepository interface {
	SendToBox(addr address.HashAddress, box int, msgID string) error
	MoveToBox(addr address.HashAddress, srcBox, dstBox int, msgID string) error

	// Message boxes
	FetchListFromBox(addr address.HashAddress, box int, since time.Time, offset, limit int) (*MessageList, error)

	// Fetch specific message contents
	FetchMessageHeader(addr address.HashAddress, box int, messageID string) (*message.Header, error)
	FetchMessageCatalog(addr address.HashAddress, box int, messageID string) ([]byte, error)
	FetchMessageBlock(addr address.HashAddress, box int, messageID, blockID string) ([]byte, error)
	FetchMessageAttachment(addr address.HashAddress, box int, messageID, attachmentID string) ([]byte, error)
}
