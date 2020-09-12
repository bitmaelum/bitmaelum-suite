package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:     "read",
	Aliases: []string{"read-message", "r"},
	Short:   "Read messages for your account",
	Long: `Read message from your account
`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()

		info := GetAccountOrDefault(vault, *rAccount)
		if info == nil {
			logrus.Fatal("No account found in vault")
			os.Exit(1)
		}

		handlers.ReadMessage(info, *rBox, *rMessageID, *rBlock)
	},
}

var (
	rAccount   *string
	rBox       *string
	rMessageID *string
	rBlock     *string
)

func init() {
	rootCmd.AddCommand(readCmd)

	rAccount = readCmd.PersistentFlags().StringP("account", "a", "", "Account")
	rBox = readCmd.PersistentFlags().StringP("box", "b", "", "Box to fetch")
	rMessageID = readCmd.PersistentFlags().StringP("message", "m", "", "Message ID")
	rBlock = readCmd.PersistentFlags().StringP("block", "", "default", "block")
}
