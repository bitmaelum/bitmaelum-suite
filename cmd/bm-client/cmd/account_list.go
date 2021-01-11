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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listAccountCmd = &cobra.Command{
	Use:   "list",
	Short: "List accounts",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		if *laccount == "" {
			displayAccounts(v)
		} else {

			addr, err := address.NewAddress(*laccount)
			if err != nil {
				fmt.Println("error: bad address: ", addr)
				os.Exit(1)
			}

			info, err := v.GetAccountInfo(*addr)
			if err != nil {
				fmt.Println("error: account not found: ", addr)
				os.Exit(1)
			}

			displayAccount(v, *info)
		}
	},
}

func short(ID string) string {
	if len(ID) < 20 {
		return ID
	}

	return ID[:8] + "..." + ID[len(ID)-8:]
}

func displayAccounts(v *vault.Vault) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Account", "Name", "Routing ID", "Server"})

	for _, a := range v.Store.Accounts {
		table.Append([]string{
			a.Address.String(),
			a.Name,
			short(a.RoutingID),
			getHostByRoutingID(a.RoutingID, "unknown host"),
		})

	}

	table.Render()
}

func displayAccount(v *vault.Vault, info vault.AccountInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetColMinWidth(1, 80)

	table.AppendBulk([][]string{
		{"Account", info.Address.String()},
		{"Account hash", info.Address.Hash().String()},
	})

	if info.Address.HasOrganisationPart() {
		table.AppendBulk([][]string{
			{"Organisation", info.Address.Org},
			{"Organisation hash", info.Address.OrgHash().String()},
		})
	}

	kp := info.GetActiveKey()
	table.AppendBulk([][]string{
		{"", ""},
		{"Full name", info.Name},
		{"Routing ID", info.RoutingID},
		{"Routing Host", getHostByRoutingID(info.RoutingID, "unknown host")},
		{"", ""},
		{"Public key", strings.Join(chunks(kp.PubKey.String(), 118), "\n")},
		{"Account mnemonic", internal.GetMnemonic(&kp.KeyPair)},
		{"Proof of work", fmt.Sprintf("%d bits", info.Pow.Bits)},
	})

	if len(info.Settings) > 0 {
		table.AppendBulk([][]string{
			{"", ""},
			{"Settings", generateSettingsTable(info.Settings)},
		})
	}

	table.Render()
}

func generateSettingsTable(settings map[string]string) string {
	ret := &strings.Builder{}
	table := tablewriter.NewWriter(ret)
	// table.SetBorder(false)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetColMinWidth(0, 30)
	table.SetColMinWidth(1, 30)

	for k, v := range settings {
		table.Append([]string{k, v})
	}

	table.Render()

	return ret.String()
}

// https://stackoverflow.com/a/61469854/1072881
func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	l := 0
	for _, r := range s {
		chunk[l] = r
		l++
		if l == chunkSize {
			chunks = append(chunks, string(chunk))
			l = 0
		}
	}
	if l > 0 {
		chunks = append(chunks, string(chunk[:l]))
	}
	return chunks
}

func getHostByRoutingID(routingID string, def string) string {
	rs := container.Instance.GetResolveService()
	info, err := rs.ResolveRouting(routingID)
	if err != nil {
		return def
	}

	return info.Routing
}

var (
	laccount *string
)

func init() {
	accountCmd.AddCommand(listAccountCmd)

	laccount = listAccountCmd.Flags().String("account", "", "Account to display")
}
