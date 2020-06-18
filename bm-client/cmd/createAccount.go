package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createAccountCmd = &cobra.Command{
	Use:   "create-account",
	Short: "Create a new account",
	Long: `Create a new account locally and upload it to a BitMaelum servrer.

This assumes you have a BitMaelum invitation token for the specific server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createAccount called")
	},
}

func init() {
	rootCmd.AddCommand(createAccountCmd)
}
