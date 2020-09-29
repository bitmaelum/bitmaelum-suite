package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/spf13/cobra"
)

var listOrganisationsCmd = &cobra.Command{
	Use:     "list-organisations",
	Aliases: []string{"list-orgs", "ls", "list"},
	Short:   "List your organisations",
	Long:    `Displays a list of all your organisations currently available`,
	Run: func(cmd *cobra.Command, args []string) {
		v := OpenVault()
		handlers.ListOrganisations(v, *orgDisplayKeys)
	},
}

var orgDisplayKeys *bool

func init() {
	rootCmd.AddCommand(listOrganisationsCmd)

	orgDisplayKeys = listOrganisationsCmd.Flags().BoolP("keys", "k", false, "Display private and public key")
}
