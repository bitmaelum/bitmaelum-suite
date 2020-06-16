package cmd

import (
	"github.com/spf13/cobra"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "Mail management",
	Long:  `Read or compose messages.`,
	Run:   SelectAndRun,
	Annotations: map[string]string{"position": "10"},
}

func init() {
	rootCmd.AddCommand(mailCmd)
}
