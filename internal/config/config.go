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

var (
	errConfigNotFound       = errors.New("cannot find config file")
	errClientConfigNotFound = errors.New("client config file not found")
	errServerConfigNotFound = errors.New("server config file not found")
	errBridgeConfigNotFound = errors.New("bridge config file not found")
)

var triedPaths []string

var (
	// ClientConfigFile Filename of client configuration
	ClientConfigFile string = "bitmaelum-client-config.yml"
	// ServerConfigFile Filename of server configuration
	ServerConfigFile string = "bitmaelum-server-config.yml"
	// BridgeConfigFile Filename of bridge configuration
	BridgeConfigFile string = "bitmaelum-bridge-config.yml"
)

// Absolute paths of the loaded configurations
var (
	LoadedClientConfigPath string
	LoadedServerConfigPath string
	LoadedBridgeConfigPath string
)

var loadedConfig string

// IsLoaded will return true when the loaded configuration is the same as the config given
func IsLoaded(config string) bool {
	return config == loadedConfig
}

// LoadClientConfig loads client configuration from given path or panic if cannot load
func LoadClientConfig(configPath string) {
	err := LoadClientConfigOrPass(configPath)
	if err != nil {
		for _, p := range triedPaths {
			logrus.Errorf("Tried path: %s", p)
		}
		logrus.Fatalf("could not load client configuration. You can generate a new configuration with: 'bm-client config init'")
	}
}

// LoadServerConfig loads server configuration from given path or panic if cannot load
func LoadServerConfig(configPath string) {
	err := LoadServerConfigOrPass(configPath)
	if err != nil {
		for _, p := range triedPaths {
			logrus.Errorf("Tried path: %s", p)
		}

		logrus.Fatalf("could not load server configuration. You can generate a new configuration with: 'bm-config init-config --server'")
	}
}

// LoadBridgeConfig loads bridge configuration from given path or panic if cannot load
func LoadBridgeConfig(configPath string) {
	err := LoadBridgeConfigOrPass(configPath)
	if err != nil {
		for _, p := range triedPaths {
			logrus.Errorf("Tried path: %s", p)
		}
		logrus.Fatalf("could not load bridge configuration. You can generate a new configuration with: 'bm-config init-config --bridge'")
	}
}

// LoadClientConfigOrPass loads client configuration, but return false if not able
func LoadClientConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		return readConfigPath(configPath, "from commandline", Client.LoadConfig, &LoadedClientConfigPath)
	}

	configPath = os.Getenv("BITMAELUM_CLIENT_CONFIG")
	if configPath != "" {
		return readConfigPath(configPath, "from BITMAELUM_CLIENT_CONFIG environment", Client.LoadConfig, &LoadedClientConfigPath)
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, ClientConfigFile)
		err = readConfigPath(p, "from hardcoded search path", Client.LoadConfig, &LoadedClientConfigPath)
		if err == nil || err != errConfigNotFound {
			return err
		}
	}

	return errClientConfigNotFound
}

// LoadServerConfigOrPass loads client configuration, but return false if not able
func LoadServerConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		return readConfigPath(configPath, "from commandline", Server.LoadConfig, &LoadedServerConfigPath)
	}

	configPath = os.Getenv("BITMAELUM_SERVER_CONFIG")
	if configPath != "" {
		return readConfigPath(configPath, "from BITMAELUM_SERVER_CONFIG environment", Server.LoadConfig, &LoadedServerConfigPath)
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, ServerConfigFile)
		err = readConfigPath(p, "from hardcoded search path", Server.LoadConfig, &LoadedServerConfigPath)

		// Set config path if this is the one found
		if err == nil {
			LoadedServerConfigPath = p
		}
		if err == nil || err != errConfigNotFound {
			return err
		}
	}

	return errServerConfigNotFound
}

// LoadBridgeConfigOrPass loads bridge configuration, but return false if not able
func LoadBridgeConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		return readConfigPath(configPath, "from commandline", Bridge.LoadConfig, &LoadedBridgeConfigPath)
	}

	configPath = os.Getenv("BITMAELUM_BRIDGE_CONFIG")
	if configPath != "" {
		return readConfigPath(configPath, "from BITMAELUM_BRIDGE_CONFIG environment", Bridge.LoadConfig, &LoadedBridgeConfigPath)
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, BridgeConfigFile)
		err = readConfigPath(p, "from hardcoded search path", Bridge.LoadConfig, &LoadedBridgeConfigPath)
		if err == nil || err != errConfigNotFound {
			return err
		}
	}

	return errBridgeConfigNotFound
}

// Expands the given path and loads the configuration
func readConfigPath(p, src string, loader func(r io.Reader) error, loadedPath *string) error {
	p, _ = homedir.Expand(p)
	p, _ = filepath.Abs(p)

	triedPaths = append(triedPaths, p+" ("+src+")")

	f, err := fs.Open(p)
	if err != nil {
		return errConfigNotFound
	}

	err = loader(f)
	if err == nil {
		*loadedPath = p
	}

	return err
}
