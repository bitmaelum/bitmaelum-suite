package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var listOrganisationsCmd = &cobra.Command{
	Use:     "list-organisations",
	Aliases: []string{"list-org", "lo"},
	Short:   "List your organisations",
	Long:    `Displays a list of all your organisations currently available`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()
		handlers.ListOrganisations(v, *orgDisplayKeys)
	},
}

var orgDisplayKeys *bool

func init() {
	rootCmd.AddCommand(listOrganisationsCmd)

	orgDisplayKeys = listOrganisationsCmd.Flags().BoolP("keys", "k", false, "Display private and public key")
}
