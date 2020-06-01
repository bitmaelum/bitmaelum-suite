package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountUpdateKeyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new private key to your account",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	accountUpdateKeyCmd.AddCommand(accountUpdateKeyAddCmd)
}
