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

package handlers

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/mailbox"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/c2h5oh/datasize"
	"github.com/olekukonko/tablewriter"
)

// AccountEntry is a single entry to display
type AccountEntry struct {
	account string   // The account of the entry
	box     string   // the box of the entry
	idx     int      // the index/position we need to display
	row     []string // actual row data for the table
}

// AccountData is the complete list of all messages we want to display.
type AccountData []AccountEntry

// Len returns the length of the account data
func (ad AccountData) Len() int {
	return len(ad)
}

// Swap will swap two account data entries in the slice
func (ad AccountData) Swap(i, j int) {
	ad[i], ad[j] = ad[j], ad[i]
}

// Less will sort account entries on account first, box second, idx third.
func (ad AccountData) Less(i, j int) bool {
	if ad[i].account < ad[j].account {
		return true
	}
	if ad[i].account > ad[j].account {
		return false
	}

	if ad[i].box < ad[j].box {
		return true
	}
	if ad[i].box > ad[j].box {
		return false
	}

	return ad[i].idx < ad[j].idx
}

// ListMessages will display message information from accounts and boxes
func ListMessages(accounts []vault.AccountInfo, since time.Time) int {
	fmt.Print("* Fetching messages from remote server(s): ")
	spinner := internal.NewSpinner(100 * time.Millisecond)
	spinner.Start()

	accountData := queryAccounts(accounts, since)

	spinner.Stop()
	fmt.Println("")

	if len(*accountData) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		headers := []string{"Account", "Box", "ID", "Subject", "From", "Date", "# Blocks", "# Attachments"}
		table.SetHeader(headers)

		for _, ad := range *accountData {
			table.Append(ad.row)
		}

		table.Render()
	}

	return len(*accountData)
}

func queryAccounts(accounts []vault.AccountInfo, since time.Time) *AccountData {
	var accountData AccountData
	var wg sync.WaitGroup

	doneCh := make(chan int)
	dataCh := make(chan AccountEntry)

	// Make a go routine for each account
	for _, info := range accounts {
		wg.Add(1)
		go func(wg *sync.WaitGroup, info vault.AccountInfo) {
			defer wg.Done()

			// Fetch routing info
			resolver := container.Instance.GetResolveService()
			routingInfo, err := resolver.ResolveRouting(info.RoutingID)
			if err != nil {
				return
			}

			client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
			if err != nil {
				return
			}

			queryBoxes(client, &info, dataCh, since)
		}(&wg, info)
	}

	// Main collector routine
	go func() {
		for {
			select {
			case data := <-dataCh:
				accountData = append(accountData, data)
			case <-doneCh:
				return
			}
		}
	}()

	// Wait until the query routines have completed
	wg.Wait()

	// Send signal to the collector routine that its done
	doneCh <- 1

	// Make sure all data is in the right order
	sort.Sort(accountData)
	return &accountData
}

func queryBoxes(client *api.API, account *vault.AccountInfo, dataCh chan AccountEntry, since time.Time) {
	mbl, err := client.GetMailboxList(account.Address.Hash())
	if err != nil {
		return
	}

	for _, mb := range mbl.Boxes {
		queryBox(client, account, fmt.Sprintf("%d", mb.ID), dataCh, since)
	}
}

func queryBox(client *api.API, account *vault.AccountInfo, box string, dataCh chan AccountEntry, since time.Time) {
	mb, err := client.GetMailboxMessages(account.Address.Hash(), box, since)
	if err != nil {
		return
	}

	// No messages in this box found
	if len(mb.Messages) == 0 {
		return
	}

	// Sort messages first
	msort := mailbox.NewMessageSort(account.GetActiveKey().PrivKey, mb.Messages, mailbox.SortDate, true)
	sort.Sort(&msort)

	// Add first entry with just the account and box number
	idx := 1
	dataCh <- AccountEntry{
		account: account.Address.String(),
		box:     box,
		idx:     idx,
		row: []string{
			account.Address.String(),
			box,
			"", "", "", "", "", "",
		},
	}

	for _, msg := range mb.Messages {
		idx++

		settings := &bmcrypto.EncryptionSettings{
			Type:          bmcrypto.CryptoType(msg.Header.Catalog.Crypto),
			TransactionID: msg.Header.Catalog.TransactionID,
		}
		key, err := bmcrypto.Decrypt(account.GetActiveKey().PrivKey, settings, msg.Header.Catalog.EncryptedKey)
		if err != nil {
			continue
		}
		catalog := &message.Catalog{}
		err = bmcrypto.CatalogDecrypt(key, msg.Catalog, catalog)
		if err != nil {
			continue
		}

		var blocks []string
		for _, b := range catalog.Blocks {
			blocks = append(blocks, b.Type)
		}

		var attachments []string
		for _, a := range catalog.Attachments {
			fs := datasize.ByteSize(a.Size)
			attachments = append(attachments, a.FileName+" ("+fs.HR()+")")
		}

		dataCh <- AccountEntry{
			account: account.Address.String(),
			box:     box,
			idx:     idx,
			row: []string{
				"",
				"",
				msg.ID,
				catalog.Subject,
				fmt.Sprintf("%s <%s>", catalog.From.Name, catalog.From.Address),
				catalog.CreatedAt.Format(time.RFC822),
				strings.Join(blocks, ","),
				strings.Join(attachments, "\n"),
			},
		}
	}
}
