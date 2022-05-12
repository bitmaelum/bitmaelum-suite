// Copyright (c) 2022 BitMaelum Authors
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

package app

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/gui/components"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/rivo/tview"
)

// BmMailAppType  is a struct that holds the current application and "global" data
type BmMailAppType struct {
	App   *tview.Application // Main TUI application
	Pages *tview.Pages       // Set of pages for the TUI application

	Vault          *vault.Vault     // Opened vault (or nil)
	CurrentAddr    *address.Address // Current selected account (or nil)
	CurrentBox     *string          // Current selected box (or nil)
	CurrentMessage *string          // Current selected message (or nil)

	StaleVault   bool // when true, the vault must be refreshed
	StaleAddr    bool // when true, the address must be refreshed
	StaleBox     bool // when true, the box must be refreshed
	StaleMessage bool // when true, the message must be refreshed

	// Cached items, can change at any time
	MailBoxLists          map[string]*api.MailboxList // Cache of all mailboxes
	CachedMailboxMessages *api.MailboxMessages        // Cache of current selected mailbox messages
	Client                *api.API                    // Current connection to the server

	// GLobal UI elements that we need to populate / use at certain occasions
	MessageAccountTree *tview.TreeView
	MessageBoxTree     *tview.TreeView
	MessageList        *components.MessageList
	MessageView        *tview.TextView
}

// MailApp is the general mail app structure
var MailApp *BmMailAppType
