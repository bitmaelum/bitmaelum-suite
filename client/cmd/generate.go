/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/jaytaph/mailv2/core/account"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new account",
	Long: `This command allows you to generate a new account`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("generating new account")

		accounts := account.LoadAccount()
		if accounts.Has(emailFlag) {
			logrus.Error("looks like this account has already be generated.")
			os.Exit(128)
		}

		account, err := accounts.GenerateAccount(emailFlag, nameFlag)
		if err != nil {
			logrus.Error("error while generating new account")
			os.Exit(128)
		}

		logrus.Infof("account for %s generated.", account.Email)
	},
}

func init() {
	accountCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVar(&emailFlag, "email", "", "Your email address to generate")
	generateCmd.Flags().StringVar(&nameFlag, "name", "", "Your name")

	_ = generateCmd.MarkFlagRequired("email")
	_ = generateCmd.MarkFlagRequired("name")
}
