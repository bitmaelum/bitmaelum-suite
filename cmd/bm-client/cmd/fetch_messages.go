package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fetchMessagesCmd = &cobra.Command{
	Use:     "fetch-messages",
	Aliases: []string{"fetch"},
	Short:   "Retrieves messages from your account(s)",
	Long:    `Connects to the BitMaelum servers and fetches new emails that are not available on your local system.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		addr, err := address.New(*fmAccount)
		if err != nil {
			logrus.Fatal(err)
		}
		info, err := vault.GetAccountInfo(*addr)
		if err != nil {
			logrus.Fatal(err)
		}

		handlers.FetchMessages(info, *fmBox, *fmCheckOnly)
	},
}

var (
	fmCheckOnly *bool
	fmAccount   *string
	fmBox       *string
)

func init() {
	rootCmd.AddCommand(fetchMessagesCmd)

	fmAccount = fetchMessagesCmd.PersistentFlags().StringP("account", "a", "", "Account")
	fmBox = fetchMessagesCmd.PersistentFlags().StringP("box", "b", "", "Box to fetch")
	fmCheckOnly = fetchMessagesCmd.Flags().Bool("check-only", false, "Check only, don't download")
}
