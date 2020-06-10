package cmd

import (
	"github.com/spf13/cobra"
)

var accountUpdateKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Key management",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: SelectAndRun,
	Annotations: map[string]string{"position": "10"},
}

func init() {
	accountUpdateCmd.AddCommand(accountUpdateKeyCmd)
}
