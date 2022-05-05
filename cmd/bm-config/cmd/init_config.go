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
	"io"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initConfigCmd represents the initConfig command
var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Creates default server and client configurations",
	Long: `Before you can run the mailserver or client, you will need a configuration file which you need to adjust 
to your own needs.

This command creates default templates that you can use as a starting point.`,
	Run: func(cmd *cobra.Command, args []string) {
		genC, _ := cmd.Flags().GetBool("client")
		genS, _ := cmd.Flags().GetBool("server")
		genB, _ := cmd.Flags().GetBool("bridge")

		// If not client or server selected, generate them both
		genAll := !genC && !genS && !genB

		if genAll || genC {
			createFile("./"+config.ClientConfigFile, config.GenerateClientConfig)
			fmt.Println("Generated client configuration file")
		}
		if genAll || genS {
			createFile("./"+config.ServerConfigFile, config.GenerateServerConfig)
			fmt.Println("Generated server configuration file")
		}
		if genAll || genB {
			createFile("./"+config.BridgeConfigFile, config.GenerateBridgeConfig)
			fmt.Println("Generated bridge configuration file")
		}
	},
}

func createFile(path string, configTemplate func(w io.Writer) error) {
	f, err := os.Create(path)
	if err != nil {
		logrus.Fatalf("Error while creating file: %v", err)
	}

	err = configTemplate(f)
	if err != nil {
		logrus.Fatalf("Error while creating file: %v", err)
	}

	err = f.Close()
	if err != nil {
		logrus.Fatalf("Error while closing file: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(initConfigCmd)

	initConfigCmd.Flags().Bool("client", false, "Generate only the client configuration")
	initConfigCmd.Flags().Bool("server", false, "Generate only the server configuration")
	initConfigCmd.Flags().Bool("bridge", false, "Generate only the bridge configuration")
}
