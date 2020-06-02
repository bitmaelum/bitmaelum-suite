package cmd

import (
	"github.com/spf13/cobra"
)

var accountUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your account",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:   SelectAndRun,
	Annotations: map[string]string{"position": "30"},
}

func init() {
	accountCmd.AddCommand(accountUpdateCmd)
}
