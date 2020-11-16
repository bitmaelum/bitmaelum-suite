// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// generateKeyCmd represents the generatekey command
var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Generates a new key",
	Long:  `This command generates a new key`,
	Run: func(cmd *cobra.Command, args []string) {
		privKey, pubKey, err := bmcrypto.GenerateKeyPair(bmcrypto.KeyType(*keyType))
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println("Private key : ", privKey.String())
		fmt.Println("Public key  : ", pubKey.String())
	},
}

var keyType *string

func init() {
	rootCmd.AddCommand(generateKeyCmd)

	keyType = generateKeyCmd.Flags().StringP("type", "t", "ed25519", "The keytype (ed25519, rsa or ecdsa)")
}
