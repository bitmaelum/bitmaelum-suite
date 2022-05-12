// Copyright (c) 2022 BitMaelum Authors
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
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initRoutingConfigCmd represents the initRoutingConfig command
var initRoutingConfigCmd = &cobra.Command{
	Use:   "generate-routing-id",
	Short: "Generates a routing ID and keypair",
	Long: `Before you can run the mail server, you will need a routing configuration file which uniquely identifies 
the server on the network.

This command creates a new routing file if one does not exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Server.Server.RoutingFile == "" {
			fmt.Println("Routing file path is not found in your server configuration.")
			os.Exit(1)
		}

		// Check if file exist
		_, err := os.Stat(config.Server.Server.RoutingFile)
		if os.IsExist(err) {
			fmt.Printf("Routing file %s already exist. I will not overwrite this file.", config.Server.Server.RoutingFile)
			os.Exit(1)
		}

		var (
			mnemonic string
			r        *config.RoutingConfig
		)

		if ok, _ := cmd.Flags().GetBool("mnemonic"); ok {
			// ask for mnemonic
			mnemonic = console.AskMnemonicPhrase()
			r, err = config.GenerateRoutingFromMnemonic(mnemonic)
			checkError(err)
		} else {
			// Generate new routing
			r, err = config.GenerateRouting()
			checkError(err)

			seed, err := hex.DecodeString(r.KeyPair.Generator)
			checkError(err)

			mnemonic, err = bmcrypto.RandomSeedToMnemonic(seed)
			checkError(err)
		}

		// Save routing
		err = config.SaveRouting(config.Server.Server.RoutingFile, r)
		checkError(err)

		logrus.Printf("Generated routing file: %s", config.Server.Server.RoutingFile)

		fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your server. If, 
for any reason, you lose this key, you will need to use the following words 
in order to recreate the key:

`)
		fmt.Print(internal.WordWrap(mnemonic, 78))
		fmt.Print(`

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.
`)
	},
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error while generating routing file: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initRoutingConfigCmd)

	initRoutingConfigCmd.Flags().Bool("mnemonic", false, "Ask for your mnemonic phrase")
}
