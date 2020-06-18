package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Compose a new message",
	Long: `This command will allow you to compose a new message and send it through your BitMaelum server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("compose called")
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
