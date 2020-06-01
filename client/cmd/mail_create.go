package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var mailCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new message box",
	Long: `Create a new message box in your account.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	mailCmd.AddCommand(mailCreateCmd)
}
