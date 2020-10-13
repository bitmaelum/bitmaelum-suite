package account

import (
	"io"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// MessageType is a simple message structure that we return as a list
type MessageType struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"header"`
	Catalog []byte         `json:"catalog"`
}

// Metatype is a structure that holds meta information about a list
type MetaType struct {
	Total    int `json:"total"`
	Returned int `json:"returned"`
	Limit    int `json:"limit"`
	Offset   int `json:"offset"`
}

// MessageList is a list of messages we return
type MessageList struct {
	Meta     MetaType
	Messages []MessageType `json:"messages"`
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
	Create(addr hash.Hash, pubKey bmcrypto.PubKey) error
	Exists(addr hash.Hash) bool
	Delete(addr hash.Hash) error
}

// KeyRepository gets and sets public keys into an account
type KeyRepository interface {
	// Public key
	StoreKey(addr hash.Hash, key bmcrypto.PubKey) error
	FetchKeys(addr hash.Hash) ([]bmcrypto.PubKey, error)
}

// OrganisationRepository gets and sets organisation settings into an account
type OrganisationRepository interface {
	StoreOrganisationSettings(addr hash.Hash, settings OrganisationSettings) error
	FetchOrganisationSettings(addr hash.Hash) (*OrganisationSettings, error)
}

// BoxRepository deals with message boxes insides an account
type BoxRepository interface {
	CreateBox(addr hash.Hash, parentBox int) error
	ExistsBox(addr hash.Hash, box int) bool
	DeleteBox(addr hash.Hash, box int) error
	GetBoxInfo(addr hash.Hash, box int) (*BoxInfo, error)
	GetAllBoxes(addr hash.Hash) ([]BoxInfo, error)
}

// MessageRepository deals with message within boxes
type MessageRepository interface {
	SendToBox(addr hash.Hash, box int, msgID string) error
	MoveToBox(addr hash.Hash, srcBox, dstBox int, msgID string) error

	// Message boxes
	FetchListFromBox(addr hash.Hash, box int, since time.Time, offset, limit int) (*MessageList, error)

	// Fetch specific message contents
	FetchMessageHeader(addr hash.Hash, box int, messageID string) (*message.Header, error)
	FetchMessageCatalog(addr hash.Hash, box int, messageID string) ([]byte, error)
	FetchMessageBlock(addr hash.Hash, box int, messageID, blockID string) ([]byte, error)
	FetchMessageAttachment(addr hash.Hash, box int, messageID, attachmentID string) (r io.ReadCloser, size int64, err error)
}
