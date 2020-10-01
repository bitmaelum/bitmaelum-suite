package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/spf13/cobra"
)

var createOrganisationInviteCmd = &cobra.Command{
	Use:   "create-organisation-invite",
	Short: "Create a new organisation invitation for a user",
	Long:  `Creates an invitation for a user for the organisation.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		handlers.CreateOrganisationInvite(vault, *orgInvOrg, *orgInvAddress, *orgInvRoutingID)
	},
}

var (
	orgInvOrg       *string
	orgInvAddress   *string
	orgInvRoutingID *string
)

func init() {
	rootCmd.AddCommand(createOrganisationInviteCmd)

	orgInvOrg = createOrganisationInviteCmd.Flags().StringP("org", "o", "", "org name")
	orgInvAddress = createOrganisationInviteCmd.Flags().StringP("addr", "a", "", "address")
	orgInvRoutingID = createOrganisationInviteCmd.Flags().StringP("routing-id", "r", "", "routing ID where this user will be invited to")

	_ = createAccountCmd.MarkFlagRequired("org")
	_ = createAccountCmd.MarkFlagRequired("org")
	_ = createAccountCmd.MarkFlagRequired("org")
}
