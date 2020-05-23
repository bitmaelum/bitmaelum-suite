/*
Copyright © 2020 Joshua Thijssen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mailv2",
	Short: "A POC for a new way of handling mail",
	Long: `This is a proof of concept on how email could function when not 
impeded by backwards compatibility. It allows us to start from scratch.`,
}

var verboseFlag *bool
var moreVerboseFlag *bool
var muchVerboseFlag *bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	if *verboseFlag == true {
		fmt.Println("INfo level")
		logger.SetLevel(logger.InfoLevel)
	}
	if *moreVerboseFlag == true {
		fmt.Println("Debug level")
		logger.SetLevel(logger.DebugLevel)
	}
	if *muchVerboseFlag == true {
		fmt.Println("trace level")
		logger.SetLevel(logger.TraceLevel)
	}
}

func init() {
	cobra.OnInitialize(initLogger)

	verboseFlag = rootCmd.PersistentFlags().BoolP("v", "", false, "verbosity level")
	moreVerboseFlag = rootCmd.PersistentFlags().BoolP("vv", "", false, "medium verbosity")
	muchVerboseFlag = rootCmd.PersistentFlags().BoolP("vvv", "", false, "highest verbosity level")
}
