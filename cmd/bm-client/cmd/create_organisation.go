package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/spf13/cobra"
)

var createOrganisationCmd = &cobra.Command{
	Use:   "create-organisation",
	Short: "Create a new organisation",
	Long: `Create a new organisation locally and upload it to the keyserver.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		handlers.CreateOrganisation(vault, *orgName, *orgValidations)
	},
}

var (
	orgName        *string
	orgValidations *[]string
)

func init() {
	rootCmd.AddCommand(createOrganisationCmd)

	orgName = createOrganisationCmd.Flags().StringP("org", "o", "", "Organisation name (...@<name>! part)")
	orgValidations = createOrganisationCmd.Flags().StringArrayP("validation", "v", nil, "validations for the organisation")

	_ = createAccountCmd.MarkFlagRequired("org")
}
