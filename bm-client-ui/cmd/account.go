package cmd

import (
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:         "account",
	Short:       "Account management",
	Long:        `Commands to manage your local accounts.`,
	Run:         SelectAndRun,
	Annotations: map[string]string{"position": "20"},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
