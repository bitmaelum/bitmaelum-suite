package cmd

import (
	"github.com/spf13/cobra"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Account resolve management",
	Long: `Manage your account resolver.`,
	Run:   SelectAndRun,
	Annotations: map[string]string{"position": "30"},
}

func init() {
	rootCmd.AddCommand(resolveCmd)
}
