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
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var organisationInviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Create a new organisation invitation for a user",
	Long:  `Creates an invitation for a user for the organisation.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()



		handlers.CreateOrganisationInvite(v, strings.TrimRight(*orgInvOrg, "!"), *orgInvAddress, *orgInvRoutingID)
	},
}

var (
	orgInvOrg       *string
	orgInvAddress   *string
	orgInvRoutingID *string
)

func init() {
	organisationCmd.AddCommand(organisationInviteCmd)

	orgInvOrg = organisationInviteCmd.Flags().StringP("organisation", "o", "", "org name")
	orgInvAddress = organisationInviteCmd.Flags().StringP("account", "a", "", "account")
	orgInvRoutingID = organisationInviteCmd.Flags().StringP("routing-id", "r", "", "routing ID where this user will be invited to")

	_ = organisationInviteCmd.MarkFlagRequired("organisation")
	_ = organisationInviteCmd.MarkFlagRequired("account")
	_ = organisationInviteCmd.MarkFlagRequired("routing-id")
}
