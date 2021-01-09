// Copyright (c) 2021 BitMaelum Authors
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
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accountKeyRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate to a new account key",
	Long: `Rotate a new key into your account. The keytype (--keytype, -k) can be either one of the following:
	
  ed25519   ED25519 elliptic curve
  rsa       RSA 2048 bit 
  rsa4096   RSA 4096 bit 

`, Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		info, err := vault.GetAccount(v, *akAccount)
		if err != nil {
			fmt.Println("cannot find account in vault")
			os.Exit(1)
		}

		kt, err := bmcrypto.FindKeyType(*keyType)
		if err != nil {
			logrus.Fatal(err)
		}

		kp, err := internal.GenerateKeypairWithRandomSeed(kt)
		if err != nil {
			logrus.Fatal(err)
		}

		info.Keys = append(info.Keys, vault.KeyPair{
			KeyPair: *kp,
			Active:  false,
		})

		err = v.Persist()
		if err != nil {
			fmt.Println("cannot save vault")
			os.Exit(1)
		}

		fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control the organisation. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:

`)
		fmt.Print(internal.WordWrap(internal.GetMnemonic(kp), 78))
		fmt.Print(`

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.

WITHOUT THESE WORDS, ALL ACCESS TO YOUR ORGANISATION IS LOST!
`)

	},
}

var (
	keyType *string
)

func init() {
	accountKeyCmd.AddCommand(accountKeyRotateCmd)

	keyType = accountKeyRotateCmd.Flags().StringP("keytype", "k", "ed25519", "The key type (defaults to ed25519)")
}
