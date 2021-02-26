// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

// MetaType is a structure that holds meta information about a list
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
	ID       int      `json:"id"`
	Total    int      `json:"total"`
	Messages []string `json:"messages"`
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
	CreateMessage(addr hash.Hash, msgID string) error
	RemoveMessage(addr hash.Hash, msgID string) error
	CopyMessage(addr hash.Hash, msgID string, boxID int) error
	MoveMessage(addr hash.Hash, msgID string, fromBoxID, toBoxID int) error
	AddToBox(addr hash.Hash, boxID int, msgID string) error
	RemoveFromBox(addr hash.Hash, boxID int, msgID string) error
	ExistsInBox(addr hash.Hash, boxID int, msgID string) bool

	// Message boxes
	FetchListFromBox(addr hash.Hash, box int, since time.Time, offset, limit int) (*MessageList, error)

	// Fetch specific message contents
	FetchMessageHeader(addr hash.Hash, messageID string) (*message.Header, error)
	FetchMessageCatalog(addr hash.Hash, messageID string) ([]byte, error)
	FetchMessageBlock(addr hash.Hash, messageID, blockID string) ([]byte, error)
	FetchMessageAttachment(addr hash.Hash, messageID, attachmentID string) (r io.ReadCloser, size int64, err error)
}
