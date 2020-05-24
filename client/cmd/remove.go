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
	"crypto/sha256"
	"encoding/base64"
	"github.com/jaytaph/mailv2/core/keys"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a key",
	Long: `Remove a key`,
	Run: func(cmd *cobra.Command, args []string) {
		hasher := sha256.New()
		hasher.Write([]byte(emailFlag))
		hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		if ! keys.HasKey(hash) {
			logger.Error("Email does not exist in the public key database")
		}
		keys.RemoveKey(hash)
	},
}

func init() {
	keysCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVar(&emailFlag, "email", "", "Email address")

	_ = removeCmd.MarkFlagRequired("email")
}
