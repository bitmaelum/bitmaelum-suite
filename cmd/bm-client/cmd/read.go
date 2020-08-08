package cmd

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages for your account",
	Long: `Read message from your account
`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := address.New(*rAccount)
		if err != nil {
			logrus.Fatal(err)
		}
		info, err := Vault.GetAccountInfo(*addr)
		if err != nil {
			logrus.Fatal(err)
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
