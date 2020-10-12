package cmd

import (
	"github.com/spf13/cobra"
)

// apiKeyCmd represents the apiKey command
var apiKeyCmd = &cobra.Command{
	Use:     "apikey",
	Short:   "Api key management",
}

func init() {
	rootCmd.AddCommand(apiKeyCmd)
}
