package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read messages for your account",
	Long: `Read message from your account
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("read called")
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
