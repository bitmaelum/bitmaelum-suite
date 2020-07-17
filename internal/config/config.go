package config

import (
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"io"
	"os"
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

// LoadClientConfig loads client configuration from given path or panic if cannot load
func LoadClientConfig(path string) {
	err := LoadClientConfigOrPass(path)
	if err != nil {
		logrus.Fatalf("cannot load client configuration: %s", err)
	}
}

// LoadServerConfig loads server configuration from given path or panic if cannot load
func LoadServerConfig(path string) {
	err := LoadServerConfigOrPass(path)
	if err != nil {
		logrus.Fatalf("cannot load server configuration: %s", err)
	}
}

// LoadClientConfigOrPass loads client configuration, but return false if not able
func LoadClientConfigOrPass(path string) error {
	var err error

	if path != "" {
		err = readConfigPath(path, Client.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	err = readConfigPath("./client-config.yml", Client.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	err = readConfigPath("~/.bitmaelum/client-config.yml", Client.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	err = readConfigPath("/etc/bitmaelum/client-config.yml", Client.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	return errors.New("cannot find client-config.yml")
}

// LoadServerConfigOrPass loads client configuration, but return false if not able
func LoadServerConfigOrPass(path string) error {
	var err error

	if path != "" {
		err = readConfigPath(path, Server.LoadConfig)
		if err == nil || err != errNotFound {
			return err
		}
	}

	err = readConfigPath("./server-config.yml", Server.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	err = readConfigPath("~/.bitmaelum/server-config.yml", Server.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	err = readConfigPath("/etc/bitmaelum/server-config.yml", Server.LoadConfig)
	if err == nil || err != errNotFound {
		return err
	}

	return errors.New("cannot find server-config.yml")
}

// Expands the given path and loads the configuration
func readConfigPath(path string, loader func(r io.Reader) error) error {
	p, _ := homedir.Expand(path)

	f, err := os.Open(p)
	if err != nil {
		return errNotFound
	}

	return loader(f)
}
