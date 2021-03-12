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
	"runtime"

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// serviceInstallCmd represents the serviceInstallCmd command
var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the service into the system",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Installing service... ")

		err := installService(getServiceNameForInstall(cmd))
		if err != nil {
			fmt.Println("ERR")
			logrus.Fatalf("Unable to install service: %v", err)
		}

		fmt.Println("OK")

		if i, _ := cmd.Flags().GetBool("start"); i {
			fmt.Print("Starting service... ")

			err = startService(getServiceName(cmd))
			if err != nil {
				fmt.Println("ERR")
				logrus.Fatalf("Unable to start service: %v", err)
			}

			fmt.Println("OK")
		}
	},
}

func installService(svc *service.Config) error {
	if svc.UserName != "" && runtime.GOOS == "windows" {
		// On windows we need the user password
		fmt.Print("\nPlease enter the password for your username \"" + svc.UserName + "\": ")
		b, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println("")

		svc.Option["Password"] = string(b)
	}

	s, err := service.New(nil, svc)
	if err != nil {
		return err
	}

	return service.Control(s, "install")
}

func init() {
	serviceInstallCmd.Flags().String("password", "", "Specify the vault password (bm-bridge) (probably not needed if running the service as a user. This is not the user password, but the vault password)")
	serviceInstallCmd.Flags().String("username", "", "Set the username to run the service as")
	serviceInstallCmd.Flags().Bool("start", false, "Start the service after install")

	serviceCmd.AddCommand(serviceInstallCmd)
}
