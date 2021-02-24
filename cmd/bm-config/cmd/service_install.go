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

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serviceInstallCmd represents the serviceInstallCmd command
var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the service into the system",
	Run: func(cmd *cobra.Command, args []string) {
		err := installService(getServiceName(cmd))
		if err != nil {
			logrus.Fatalf("Unable to install service: %v", err)
		}

		fmt.Println("Service installed")
	},
}

func installService(svc *service.Config) error {
	s, err := service.New(nil, svc)
	if err != nil {
		return err
	}

	return service.Control(s, "install")
}

func init() {
	serviceInstallCmd.Flags().String("imaphost", "", "Set the host for the IMAP server (bm-bridge)")
	serviceInstallCmd.Flags().String("smtphost", "", "Set the host for the SMTP server (bm-bridge)")
	serviceInstallCmd.Flags().String("password", "", "Specify the vault password (usually not needed)")
	serviceInstallCmd.Flags().Bool("bm-server", false, "Manage bm-server service")
	serviceInstallCmd.Flags().Bool("bm-bridge", false, "Manage bm-bridge service")

	serviceCmd.AddCommand(serviceInstallCmd)
}
