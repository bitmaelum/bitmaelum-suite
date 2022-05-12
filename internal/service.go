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

package internal

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

const (
	windowsOS = "windows"
	linuxOS   = "linux"
	macOS     = "darwin"
)

type options struct {
	ImapHost       string `long:"imaphost" description:"Host:Port to imap server from" required:"false"`
	SMTPHost       string `long:"smtphost" description:"Host:Port to smtp server from" required:"false"`
	Config         string `short:"c" long:"config" description:"Path to your configuration file"`
	Password       string `short:"p" long:"password" description:"Vault password" default:""`
	UserName       string `long:"username" description:"Username to run the service as" default:""`
	GatewayAccount string `long:"gatewayaccount" description:"Account to use to check for pending outgoing mails" required:"false"`
}

var opts options

// GetBMServerService will return the service info
func GetBMServerService(executable string) *service.Config {
	ParseOptions(&opts)

	var arguments []string
	arguments = append(arguments, "--service")

	config.LoadServerConfig(opts.Config)
	arguments = append(arguments, "--config="+config.LoadedServerConfigPath)

	svcConfig := getServiceConfig("BM-Server", "BitMaelum server", "BitMaelum server service", executable, arguments)

	if svcConfig == nil {
		return nil
	}

	if opts.UserName != "" {
		svcConfig.UserName = opts.UserName
	}

	return svcConfig
}

// GetBMBridgeService will return the service info
func GetBMBridgeService(executable string) *service.Config {
	var arguments []string

	if executable != "" {
		// install mode
		ParseOptions(&opts)

		arguments = append(arguments, "--service")

		if opts.Password != "" {
			arguments = append(arguments, "--password="+opts.Password)
		}

		config.LoadBridgeConfig(opts.Config)
		arguments = append(arguments, "--config="+config.LoadedBridgeConfigPath)
	}

	// Get current user
	user, err := user.Current()
	if err != nil {
		return nil
	}

	svcConfig := getServiceConfig("BM-Bridge", "BitMaelum email bridge", "BitMaelum email bridge service", executable, arguments)

	if svcConfig == nil {
		return nil
	}

	if opts.UserName != "" {
		svcConfig.UserName = opts.UserName
	} else {
		svcConfig.UserName = user.Username
	}

	return svcConfig
}

func getServiceConfig(name, displayName, description, executable string, arguments []string) *service.Config {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"

	svcConfig := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Arguments:   arguments,
		Option:      options,
	}

	switch runtime.GOOS {
	case linuxOS:
		svcConfig.Dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"}

	case windowsOS:
		svcConfig.Dependencies = []string{"Tcpip"}
	}

	if executable != "" {
		if runtime.GOOS == windowsOS {
			executable = executable + ".exe"
		}
		svcConfig.Executable = executable

		if _, err := os.Stat(svcConfig.Executable); os.IsNotExist(err) {
			// file not found, try on path
			svcConfig.Executable, err = exec.LookPath(svcConfig.Executable)
			if err != nil {
				// file not found, try on the same path where we are
				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
				svcConfig.Executable = filepath.Join(dir, executable)
				if _, err := os.Stat(svcConfig.Executable); os.IsNotExist(err) {
					logrus.Fatalf("Please run this command on the same path where \"%s\" is located or add it to your path", executable)
					return nil
				}
			}
		}

	}

	return svcConfig
}
