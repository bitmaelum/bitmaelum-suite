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
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/cobra"
)

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new config file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		_, ok := os.Stat(*cFile)
		if ok == nil {
			fmt.Println("Configuration file already exist. Not overwriting the current one.")
			os.Exit(1)
		}

		err := createConfigFile(*cFile)
		if err != nil {
			fmt.Println("Error while creating file: ", err)
			os.Exit(1)
		}

		fmt.Printf("successfully created a new configuration file at %s\n", *cFile)
	},
	Annotations: map[string]string{
		"dont_load_config": "true",
	},
}

func createConfigFile(p string) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}

	err = config.GenerateClientConfig(f)
	if err != nil {
		return err
	}

	return f.Close()
}

var (
	cFile *string
)

func init() {
	configCmd.AddCommand(configInitCmd)

	cFile = configInitCmd.Flags().StringP("file", "f", "./bitmaelum-client-config.yml", "Path to configuration file")
}
