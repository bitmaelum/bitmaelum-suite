package cmd

import (
	"github.com/spf13/cobra"
)

var createOrganisationCmd = &cobra.Command{
	Use:   "create-organisation",
	Short: "Create a new organisation",
	Long: `Create a new organisation locally and upload it to the keyserver.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		// vault := OpenVault()

		//handlers.CreateOrganisation(vault, *orgName, *orgDefaultRouting)
	},
}

var orgName, orgDefaultRouting *string

func init() {
	rootCmd.AddCommand(createOrganisationCmd)

	orgName = createOrganisationCmd.Flags().String("org", "", "Organisation name (...@<name>! part)")
	orgDefaultRouting = createOrganisationCmd.Flags().String("routing", "", "Default routing info for the organisation")

	_ = createAccountCmd.MarkFlagRequired("org")
	_ = createAccountCmd.MarkFlagRequired("routing")
}
