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
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"github.com/jaytaph/mailv2/core/keys"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
)


// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new key",
	Long: `Add a new public key for an email address`,
	Run: func(cmd *cobra.Command, args []string) {
		sum := sha256.Sum256([]byte(emailFlag))
		hash := hex.EncodeToString(sum[:])

		if keys.HasKey(hash) {
			logger.Error("Email already exist in the public key database")
		}

		var reader io.Reader

		if keyFileFlag == "-" {
			reader = os.Stdin
		} else {
			f, err := os.Open(keyFileFlag)
			if err != nil {
				log.Panicf("cannot open file %s: %s", keyFileFlag, err)
			}

			reader = f
			defer f.Close()
		}

		pubKey, err := ioutil.ReadAll(reader)
		if err != nil {
			logger.Panicf("cannot read file: %s", err)
		}

		// Sanity check to see if our file really contains a public key
		block, _ := pem.Decode(pubKey)
		_, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Panicf("file doesn't seem to contain a public key: %s", err)
		}

		keys.AddKey(hash, string(pubKey))
	},
}

func init() {
	keysCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&emailFlag, "email", "", "Email address")
	addCmd.Flags().StringVar(&keyFileFlag, "keyfile", "", "path to key or use - for reading from STDIN")

	_ = addCmd.MarkFlagRequired("email")
	_ = addCmd.MarkFlagRequired("keyfile")
}
