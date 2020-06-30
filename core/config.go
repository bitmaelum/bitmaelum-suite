package core

import (
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/mitchellh/go-homedir"
	"log"
)

/**
 * Configuration is found the following way:
 *
 * 1. Check path is not empty and config found in path
 * 2. Check for config (*-config.yml) in current directory
 * 3. Check for config (*-config.yml) in directory ~/.bitmaelum
 * 4. Check for config (*-config.yml) in directory /etc/bitmaelum
 * 5. Error (or pass)
 */

// LoadClientConfig loads client configuration from given path or panic if cannot load
func LoadClientConfig(path string) {
	loaded := LoadClientConfigOrPass(path)
	if !loaded {
		log.Fatalf("cannot load client configuration")
	}
}

// LoadServerConfig loads server configuration from given path or panic if cannot load
func LoadServerConfig(path string) {
	loaded := LoadServerConfigOrPass(path)
	if !loaded {
		log.Fatal("cannot load server configuration")
	}
}

// LoadClientConfigOrPass loads client configuration, but return false if not able
func LoadClientConfigOrPass(path string) bool {
	var err error

	if path != "" {
		err = readConfigPath(path, config.Client.LoadConfig)
		if err == nil {
			return true
		}
	}

	err = readConfigPath("~/.bitmaelum/client-config.yml", config.Client.LoadConfig)
	if err == nil {
		return true
	}

	err = readConfigPath("/etc/bitmaelum/client-config.yml", config.Client.LoadConfig)
	if err == nil {
		return true
	}

	return false
}

// LoadServerConfigOrPass loads client configuration, but return false if not able
func LoadServerConfigOrPass(path string) bool {
	var err error

	if path != "" {
		err = readConfigPath(path, config.Server.LoadConfig)
		if err == nil {
			return true
		}
	}

	err = readConfigPath("~/.bitmaelum/server-config.yml", config.Server.LoadConfig)
	if err == nil {
		return true
	}

	err = readConfigPath("/etc/bitmaelum/server-config.yml", config.Server.LoadConfig)
	if err == nil {
		return true
	}

	return false
}

// Expands the given path and loads the configuration
func readConfigPath(path string, loader func(string) error) error {
	p, _ := homedir.Expand(path)
	return loader(p)
}
