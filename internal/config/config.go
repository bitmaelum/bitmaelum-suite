package config

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var errNotFound = errors.New("cannot find config file")

/**
 * Configuration is found the following way:
 *
 * 1. Check path is not empty and config found in path
 * 2. Check for config (*-config.yml) in current directory
 * 3. Check for config (*-config.yml) in directory ~/.bitmaelum
 * 4. Check for config (*-config.yml) in directory /etc/bitmaelum
 * 5. Error (or pass)
 *
 * Config assumes that all paths are expanded with homedir.Expand
 */

var triedPaths []string

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
		if err == nil || err != errNotFound {
			return err
		}
	}

	configPath = os.Getenv("BITMAELUM_CLIENT_CONFIG")
	if configPath != "" {
		err = readConfigPath(configPath, Server.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	// try on our search paths
	for _, p := range getSearchPaths() {
		p = filepath.Join(p, "client-config.yml")
		err = readConfigPath(p, Client.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	return errors.New("cannot find client-config.yml")
}

// LoadServerConfigOrPass loads client configuration, but return false if not able
func LoadServerConfigOrPass(configPath string) error {
	var err error

	// Try custom path first
	if configPath != "" {
		err = readConfigPath(configPath, Server.LoadConfig)
		if err == nil || err != errNotFound {
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
		p = filepath.Join(p, "server-config.yml")
		err = readConfigPath(p, Server.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	return errors.New("cannot find server-config.yml")
}

// Expands the given path and loads the configuration
func readConfigPath(p string, loader func(r io.Reader) error) error {
	p, _ = homedir.Expand(p)

	triedPaths = append(triedPaths, p)

	f, err := os.Open(p)
	if err != nil {
		return errNotFound
	}

	return loader(f)
}
