package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var createOrganisationCmd = &cobra.Command{
	Use:   "create-organisation",
	Short: "Create a new organisation",
	Long: `Create a new organisation locally and upload it to the keyserver.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		handlers.CreateOrganisation(v, *orgAddr, *orgFullName, *orgValidations)
	},
}

var (
	orgAddr        *string
	orgFullName    *string
	orgValidations *[]string
)

func init() {
	rootCmd.AddCommand(createOrganisationCmd)

	orgAddr = createOrganisationCmd.Flags().StringP("org", "o", "", "Organisation address (...@<name>! part)")
	orgFullName = createOrganisationCmd.Flags().StringP("name", "n", "", "Actual name (Acme Inc.)")
	orgValidations = createOrganisationCmd.Flags().StringArray("val", nil, "validations for the organisation")

	_ = createAccountCmd.MarkFlagRequired("org")
}
