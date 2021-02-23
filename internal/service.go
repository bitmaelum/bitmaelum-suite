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

package internal

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/kardianos/service"
)

const (
	windowsOS = "windows"
)

// GetBMServerService will return the service info
func GetBMServerService(executable string) *service.Config {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"

	svcConfig := &service.Config{
		Name:        "BM-Server",
		DisplayName: "BitMaelum server",
		Description: "BitMaelum server service",
		Arguments:   []string{"--service"},
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	if executable != "" {
		svcConfig.Executable = executable
		if runtime.GOOS == windowsOS {
			svcConfig.Executable = executable + ".exe"
		}

		if _, err := os.Stat(svcConfig.Executable); os.IsNotExist(err) {
			svcConfig.Executable, err = exec.LookPath(svcConfig.Executable)
			if err != nil {
				return nil
			}
		}

	}

	return svcConfig
}

// GetBMBridgeService will return the service info
func GetBMBridgeService(executable, imaphost, smtphost string) *service.Config {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"

	if imaphost == "" {
		imaphost = "127.0.0.1:1143"
	}

	if smtphost == "" {
		smtphost = "127.0.0.1:1025"
	}

	svcConfig := &service.Config{
		Name:        "BM-Bridge",
		DisplayName: "BitMaelum email bridge",
		Description: "BitMaelum email bridge service",
		Arguments:   []string{"--service", "--imaphost=" + imaphost, "--smtphost=" + smtphost},
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	if executable != "" {
		svcConfig.Executable = executable
		if runtime.GOOS == windowsOS {
			svcConfig.Executable = executable + ".exe"
		}

		if _, err := os.Stat(svcConfig.Executable); os.IsNotExist(err) {
			svcConfig.Executable, err = exec.LookPath(svcConfig.Executable)
			if err != nil {
				return nil
			}
		}

	}

	return svcConfig
}
