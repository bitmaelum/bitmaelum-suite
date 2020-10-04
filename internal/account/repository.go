package account

import (
	"io"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
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

// OrganisationSettings defines settings for organisations
type OrganisationSettings struct {
	OnlyAllowAccountsOnMainServer bool `json:"only_allow_main_server_accounts"`
}

// BoxInfo returns information about the given message box
type BoxInfo struct {
	ID    int `json:"id"`
	Total int `json:"total"`
}

// Repository is the main repository that needs to be implemented. It's pretty big
type Repository interface {
	AddressRepository
	KeyRepository
	BoxRepository
	MessageRepository
	OrganisationRepository
}

// AddressRepository creates, checks or deletes complete accounts. Address is not the correct word for this.
type AddressRepository interface {
	Create(addr address.Hash, pubKey bmcrypto.PubKey) error
	Exists(addr address.Hash) bool
	Delete(addr address.Hash) error
}

// KeyRepository gets and sets public keys into an account
type KeyRepository interface {
	// Public key
	StoreKey(addr address.Hash, key bmcrypto.PubKey) error
	FetchKeys(addr address.Hash) ([]bmcrypto.PubKey, error)
}

// OrganisationRepository gets and sets organisation settings into an account
type OrganisationRepository interface {
	StoreOrganisationSettings(addr address.Hash, settings OrganisationSettings) error
	FetchOrganisationSettings(addr address.Hash) (*OrganisationSettings, error)
}

// BoxRepository deals with message boxes insides an account
type BoxRepository interface {
	CreateBox(addr address.Hash, parentBox int) error
	ExistsBox(addr address.Hash, box int) bool
	DeleteBox(addr address.Hash, box int) error
	GetBoxInfo(addr address.Hash, box int) (*BoxInfo, error)
	GetAllBoxes(addr address.Hash) ([]BoxInfo, error)
}

// MessageRepository deals with message within boxes
type MessageRepository interface {
	SendToBox(addr address.Hash, box int, msgID string) error
	MoveToBox(addr address.Hash, srcBox, dstBox int, msgID string) error

	// Message boxes
	FetchListFromBox(addr address.Hash, box int, since time.Time, offset, limit int) (*MessageList, error)

	// Fetch specific message contents
	FetchMessageHeader(addr address.Hash, box int, messageID string) (*message.Header, error)
	FetchMessageCatalog(addr address.Hash, box int, messageID string) ([]byte, error)
	FetchMessageBlock(addr address.Hash, box int, messageID, blockID string) ([]byte, error)
	FetchMessageAttachment(addr address.Hash, box int, messageID, attachmentID string) (r io.ReadCloser, size int64, err error)
}
