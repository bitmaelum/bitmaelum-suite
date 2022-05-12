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

package cmd

import (
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	addressColName      = "Address"
	hashColName         = "Hash"
	publicKeyColName    = "Public Key"
	organisationColName = "Organisation"
	proofOfWorkColName  = "Proof of work"
	routingIDColName    = "Routing ID"
	routingColName      = "Routing"
)

var resolverQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "query resolver",
	Run: func(cmd *cobra.Command, args []string) {

		if *raAccount != "" {
			queryAccount(*raAccount)
		}

		if *raOrganisation != "" {
			queryOrganisation(*raOrganisation)
		}

		if *raRouting != "" {
			queryRouting(*raRouting)
		}
	},
}

func queryAccount(account string) {
	addr, err := address.NewAddress(account)
	if err != nil {
		fatal("bad address: ", addr)
	}

	rs := container.Instance.GetResolveService()
	info, err := rs.ResolveAddress(addr.Hash())
	if err != nil {
		fatal("cannot query account: ", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)

	table.AppendBulk([][]string{
		{addressColName, account},
		{hashColName, info.Hash},
		{publicKeyColName, strings.Join(chunks(info.PublicKey.String(), 50), "\n")},
		{proofOfWorkColName, info.Pow},
		{"", ""},
		{routingIDColName, info.RoutingID},
		{routingColName, info.RoutingInfo.Routing},
	})

	table.Render()
}

func queryOrganisation(organisation string) {
	orgHash := hash.New(organisation)

	rs := container.Instance.GetResolveService()
	info, err := rs.ResolveOrganisation(orgHash)
	if err != nil {
		fatal("cannot query organisation: ", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)

	table.AppendBulk([][]string{
		{organisationColName, organisation},
		{hashColName, info.Hash},
		{publicKeyColName, strings.Join(chunks(info.PublicKey.String(), 50), "\n")},
		{proofOfWorkColName, info.Pow},
		{"", ""},
	})

	if len(info.Validations) > 0 {
		table.AppendBulk([][]string{
			{"", ""},
			{"Validations", toMultiLine(info.Validations)},
		})
	}

	table.Render()
}

func toMultiLine(validations []organisation.ValidationType) string {
	var ret []string

	for _, val := range validations {
		ret = append(ret, val.String())
	}

	return strings.Join(ret, "\n")
}

func queryRouting(routingID string) {
	rs := container.Instance.GetResolveService()
	info, err := rs.ResolveRouting(routingID)
	if err != nil {
		fatal("cannot query routing: ", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("|")
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)

	table.AppendBulk([][]string{
		{hashColName, info.Hash},
		{publicKeyColName, strings.Join(chunks(info.PublicKey.String(), 50), "\n")},
		{routingColName, info.Routing},
	})

	table.Render()
}

var (
	raAccount      *string
	raOrganisation *string
	raRouting      *string
)

func init() {
	resolverCmd.AddCommand(resolverQueryCmd)

	raAccount = resolverQueryCmd.Flags().StringP("account", "a", "", "Account to query")
	raOrganisation = resolverQueryCmd.Flags().StringP("organisation", "o", "", "Organisation to query")
	raRouting = resolverQueryCmd.Flags().StringP("routing", "r", "", "Routing to query")
}
