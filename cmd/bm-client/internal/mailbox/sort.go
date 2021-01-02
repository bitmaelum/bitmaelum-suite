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

package mailbox

import (
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// MessageSort is a generic message sorter
type MessageSort struct {
	Messages       []api.MailboxMessagesMessage
	key            *bmcrypto.PrivKey
	openedCatalogs map[string]*message.Catalog
	sorter         func(c1, c2 *message.Catalog, asc bool) bool
	asc            bool
}

// SortField defines which kind of sorting we need to do
type SortField int

// Specific sorting fields
const (
	SortDate SortField = iota
	SortFrom
	SortSubject
)

// NewMessageSort will create a new sorted based on the given sortfield and ascending/descending order
func NewMessageSort(key *bmcrypto.PrivKey, messages []api.MailboxMessagesMessage, field SortField, asc bool) MessageSort {
	ms := MessageSort{
		Messages:       messages,
		key:            key,
		openedCatalogs: make(map[string]*message.Catalog),
		asc:            asc,
	}

	switch field {
	case SortDate:
		ms.sorter = func(c1, c2 *message.Catalog, asc bool) bool {
			if asc {
				return c1.CreatedAt.After(c2.CreatedAt)
			}
			return c1.CreatedAt.Before(c2.CreatedAt)
		}
	case SortFrom:
		ms.sorter = func(c1, c2 *message.Catalog, asc bool) bool {
			if asc {
				return strings.ToLower(c1.From.Name) > strings.ToLower(c2.From.Name)
			}
			return strings.ToLower(c1.From.Name) < strings.ToLower(c2.From.Name)
		}
	case SortSubject:
		ms.sorter = func(c1, c2 *message.Catalog, asc bool) bool {
			if asc {
				return strings.ToLower(c1.Subject) > strings.ToLower(c2.Subject)
			}
			return strings.ToLower(c1.Subject) < strings.ToLower(c2.Subject)
		}
	}

	return ms
}

// Len returns the length of the message slice
func (ms *MessageSort) Len() int {
	return len(ms.Messages)
}

// Swap will swap two messages in the slice
func (ms *MessageSort) Swap(i, j int) {
	ms.Messages[i], ms.Messages[j] = ms.Messages[j], ms.Messages[i]
}

// Less will return if message i is less then message j through the sorter function
func (ms *MessageSort) Less(i, j int) bool {
	c1 := ms.openOrGetCatalog(ms.Messages[i])
	c2 := ms.openOrGetCatalog(ms.Messages[j])

	return ms.sorter(c1, c2, ms.asc)
}

// openOrGetCatalog will either open and decrypt a catalog, or opens a cached catalog
func (ms *MessageSort) openOrGetCatalog(msg api.MailboxMessagesMessage) *message.Catalog {
	cat, ok := ms.openedCatalogs[msg.ID]
	if ok {
		return cat
	}

	key, _ := bmcrypto.Decrypt(*ms.key, msg.Header.Catalog.TransactionID, msg.Header.Catalog.EncryptedKey)
	cat = &message.Catalog{}
	_ = bmcrypto.CatalogDecrypt(key, msg.Catalog, cat)

	ms.openedCatalogs[msg.ID] = cat
	return cat
}
