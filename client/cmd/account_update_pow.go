package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var accountUpdatePowCmd = &cobra.Command{
	Use:   "pow",
	Short: "Update your proof-of-work",
	Long: `Proof-of-work proofs that your account has done work. This is needed in order 
to combat fake accounts as it gets too expensive to create large volumes of fake accounts.
Message servers can demand a certain work to be done.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pow called")
	},
}

func init() {
	accountUpdateCmd.AddCommand(accountUpdatePowCmd)
}
