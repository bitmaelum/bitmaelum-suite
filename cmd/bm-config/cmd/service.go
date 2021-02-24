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
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// serviceCmd represents the serviceCmd command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Service management",
}

func getServiceName(cmd *cobra.Command) *service.Config {
	if i, _ := cmd.InheritedFlags().GetBool("bm-server"); i {
		return internal.GetBMServerService("")
	}

	if i, _ := cmd.InheritedFlags().GetBool("bm-bridge"); i {
		return internal.GetBMBridgeService("")
	}

	cmd.Help()
	os.Exit(0)
	return nil
}

func getServiceNameForInstall(cmd *cobra.Command) *service.Config {
	if i, _ := cmd.InheritedFlags().GetBool("bm-server"); i {
		return internal.GetBMServerService("bm-server")
	}

	if i, _ := cmd.InheritedFlags().GetBool("bm-bridge"); i {
		return internal.GetBMBridgeService("bm-bridge")
	}

	cmd.Help()
	os.Exit(0)
	return nil
}

func init() {
	serviceCmd.PersistentFlags().Bool("bm-server", false, "Manage bm-server service")
	serviceCmd.PersistentFlags().Bool("bm-bridge", false, "Manage bm-bridge service")

	rootCmd.AddCommand(serviceCmd)
}
