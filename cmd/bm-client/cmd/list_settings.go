package cmd

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listSettingsCmd = &cobra.Command{
	Use:   "list-settings",
	Short: "List settings for your account",
	Long:  `Displays a list of all your settings`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		info := vault.GetAccountOrDefault(v, *lsAddr)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Key", "Value"})

		table.Append([]string{"Name", info.Name})

		if info.Settings != nil {
			for k, v := range info.Settings {
				table.Append([]string{k, v})
			}
		}

		table.Render()
	},
}

var lsAddr *string

func init() {
	rootCmd.AddCommand(listSettingsCmd)

	lsAddr = listSettingsCmd.Flags().String("address", "", "Address to display settings")
}
