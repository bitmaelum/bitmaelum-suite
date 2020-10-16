// Copyright (c) 2020 BitMaelum Authors
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

package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

var errNotFound = errors.New("cannot find config file")

var triedPaths []string

var (
	// ClientConfigFile Filename of client configuration
	ClientConfigFile string = "client-config.yml"
	// ServerConfigFile Filename of server configuration
	ServerConfigFile string = "server-config.yml"
)

// LoadClientConfig loads client configuration from given path or panic if cannot load
func LoadClientConfig(configPath string) {
	err := LoadClientConfigOrPass(configPath)
	if err != nil {
		for _, p := range triedPaths {
			logrus.Errorf("Tried path: %s", p)
		}
		logrus.Fatalf("could not load client configuration")
	}
}

// LoadServerConfig loads server configuration from given path or panic if cannot load
func LoadServerConfig(configPath string) {
	err := LoadServerConfigOrPass(configPath)
	if err != nil {
		for _, p := range triedPaths {
			logrus.Errorf("Tried path: %s", p)
		}

		logrus.Fatalf("could not load server configuration")
	}
}

// LoadClientConfigOrPass loads client configuration, but return false if not able
func LoadClientConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		err = readConfigPath(configPath, Client.LoadConfig)
		if err != nil {
			return err
		}
	}

	configPath = os.Getenv("BITMAELUM_CLIENT_CONFIG")
	if configPath != "" {
		err = readConfigPath(configPath, Client.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, ClientConfigFile)
		err = readConfigPath(p, Client.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	return errors.New("cannot find " + ClientConfigFile)
}

// LoadServerConfigOrPass loads client configuration, but return false if not able
func LoadServerConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		err = readConfigPath(configPath, Server.LoadConfig)
		if err != nil {
			return err
		}
	}

	configPath = os.Getenv("BITMAELUM_SERVER_CONFIG")
	if configPath != "" {
		err = readConfigPath(configPath, Server.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, ServerConfigFile)
		err = readConfigPath(p, Server.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	return errors.New("cannot find " + ServerConfigFile)
}

// Expands the given path and loads the configuration
func readConfigPath(p string, loader func(r io.Reader) error) error {
	p, _ = homedir.Expand(p)

	triedPaths = append(triedPaths, p)

	f, err := fs.Open(p)
	if err != nil {
		return errNotFound
	}

	return loader(f)
}
