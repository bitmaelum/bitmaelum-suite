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

package layout

import (
	"fmt"
	"sort"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-mail/app"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Refresh() {

	// Resync the vault if needed
	if app.MailApp.Vault != nil && app.MailApp.StaleVault {
		RefreshAccountList()
		RefreshAccountScreen()

		app.MailApp.StaleVault = false
	}

	if app.MailApp.StaleAddr {
		// Check if the have cached the MBL list, if not, load it
		_, ok := app.MailApp.MailBoxLists[app.MailApp.CurrentAddr.String()]
		if !ok {
			mbl, err := getMailboxList()
			if err != nil {
				return
			}
			app.MailApp.MailBoxLists[app.MailApp.CurrentAddr.String()] = mbl
		}

		// Update all things which uses the address
		RefreshBoxes()

		app.MailApp.StaleAddr = false
	}

	if app.MailApp.StaleBox {
		// Load message from box
		var err error
		app.MailApp.CachedMailboxMessages, err = getMessagesFromBox(app.MailApp.CurrentAddr, app.MailApp.CurrentBox)
		if err != nil {
			return
		}

		// Update all things which uses the box
		RefreshMessageList()
		app.MailApp.StaleBox = false
	}

	if app.MailApp.StaleMessage {
		// Reload message info
		// Update all things which uses the message
		app.MailApp.StaleMessage = false
	}
}

func getMessagesFromBox(addr *address.Address, box *string) (*api.MailboxMessages, error) {
	info, err := authenticate()
	if err != nil {
		return nil, err
	}

	return app.MailApp.Client.GetMailboxMessages(info.Address.Hash(), *app.MailApp.CurrentBox, time.Time{})
}

func getMailboxList() (*api.MailboxList, error) {
	info, err := authenticate()
	if err != nil {
		return nil, err
	}

	return app.MailApp.Client.GetMailboxList(info.Address.Hash())
}

func authenticate() (*vault.AccountInfo, error) {
	info, err := app.MailApp.Vault.GetAccountInfo(*app.MailApp.CurrentAddr)
	if err != nil {
		return nil, err
	}

	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, err
	}

	app.MailApp.Client, err = api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, nil)
	if err != nil {
		return nil, err
	}
	
	return info, nil
}

func RefreshBoxes() {
	if app.MailApp.CurrentAddr == nil {
		return
	}

	addr := app.MailApp.CurrentAddr.String()

	tree := tview.NewTreeNode(addr).SetColor(tcell.ColorYellow).SetExpanded(true).SetSelectable(false)
	for _, box := range app.MailApp.MailBoxLists[addr].Boxes {
		if box.Total == 0 {
			tree.AddChild(tview.NewTreeNode(fmt.Sprintf("box-%d", box.ID)))
		} else {
			tree.AddChild(tview.NewTreeNode(fmt.Sprintf("box-%d (%d)", box.ID, box.Total)))
		}
	}

	// Select first entry
	app.MailApp.MessageBoxTree.SetCurrentNode(tree.GetChildren()[0])
	app.MailApp.MessageBoxTree.SetRoot(tree)
	app.MailApp.App.SetFocus(app.MailApp.MessageBoxTree)
}

type sortedAccounts []vault.AccountInfo

func (s sortedAccounts) Len() int {
	return len(s)
}

func (s sortedAccounts) Less(i, j int) bool {
	return s[i].Address.String() < s[j].Address.String()
}

func (s sortedAccounts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func RefreshAccountList() {
	if app.MailApp.Vault == nil {
		return
	}

	tree := tview.NewTreeNode("accounts").SetColor(tcell.ColorYellow).SetExpanded(true).SetSelectable(false)

	sort.Sort(sortedAccounts(app.MailApp.Vault.Store.Accounts))
	for _, acc := range app.MailApp.Vault.Store.Accounts {
		tree.AddChild(tview.NewTreeNode(acc.Address.String()).SetColor(tcell.ColorYellow))
	}

	// Select first entry
	app.MailApp.MessageAccountTree.SetCurrentNode(tree.GetChildren()[0])
	app.MailApp.MessageAccountTree.SetRoot(tree)
	app.MailApp.App.SetFocus(app.MailApp.MessageAccountTree)
}

func RefreshAccountScreen() {
	if app.MailApp.Vault == nil {
		AccountList.Clear()
		OrganisationList.Clear()
		return
	}

	for _, acc := range app.MailApp.Vault.Store.Accounts {
		AccountList.AddItem(acc.Name+" <"+acc.Address.String()+">", "", rune(0), nil)
	}
	for _, org := range app.MailApp.Vault.Store.Organisations {
		OrganisationList.AddItem(org.FullName+" <...@"+org.Addr+">", "", rune(0), nil)
	}
}

func RefreshMessageList() {
	info, err := app.MailApp.Vault.GetAccountInfo(*app.MailApp.CurrentAddr)
	if err != nil {
		return
	}

	var msgs []message.DecryptedMessage
	for _, msg := range app.MailApp.CachedMailboxMessages.Messages {
		em := message.EncryptedMessage{
			ID:                       msg.ID,
			Header:                   &msg.Header,
			Catalog:                  msg.Catalog,
			GenerateBlockReader:      app.MailApp.Client.GenerateAPIBlockReader(info.Address.Hash()),
			GenerateAttachmentReader: app.MailApp.Client.GenerateAPIAttachmentReader(info.Address.Hash()),
		}

		// decrypt message
		keys := info.Keys
		key, err := info.FindKey(msg.Header.To.Fingerprint)
		if err == nil {
			keys = []vault.KeyPair{*key}
		}

		// Iterate all keys (or the single key), and see if we can decrypt the message
		var decryptedMsg *message.DecryptedMessage
		for _, key := range keys {
			var err error
			decryptedMsg, err = em.Decrypt(key.PrivKey)
			if err == nil {
				break
			}
			decryptedMsg = nil
		}

		if decryptedMsg != nil {
			msgs = append(msgs, *decryptedMsg)
		}
	}

	app.MailApp.MessageList.Items = msgs

	app.MailApp.App.SetFocus(app.MailApp.MessageList)
}
