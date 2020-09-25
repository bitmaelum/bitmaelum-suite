package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var fetchMessagesCmd = &cobra.Command{
	Use:     "fetch-messages",
	Aliases: []string{"fetch"},
	Short:   "Retrieves messages from your account(s)",
	Long:    `Connects to the BitMaelum servers and fetches new emails that are not available on your local system.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		info := GetAccountOrDefault(vault, *fmAccount)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		// Fetch routing info
		resolver := container.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("Cannot find routing ID for this account")
			os.Exit(1)
		}

		handlers.FetchMessages(info, routingInfo, *fmBox)
	},
}

var (
	// fmCheckOnly *bool
	fmAccount *string
	fmBox     *string
)

func init() {
	rootCmd.AddCommand(fetchMessagesCmd)

	fmAccount = fetchMessagesCmd.PersistentFlags().StringP("account", "a", "", "Account")
	fmBox = fetchMessagesCmd.PersistentFlags().StringP("box", "b", "", "Box to fetch")
	// fmCheckOnly = fetchMessagesCmd.Flags().Bool("check-only", false, "Check only, don't download")
}
