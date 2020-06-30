package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Client keeps all client configuration settings
var Client ClientConfig = ClientConfig{}

// Basically, our config is inside the "config" section. So we load the whole file and only store the Cfg section
type wrappedClientConfig struct {
	Cfg ClientConfig `yaml:"config"`
}

// ClientConfig is the representation of the client configuration
type ClientConfig struct {
	Accounts struct {
		Path        string `yaml:"path"`
		ProofOfWork int    `yaml:"proof_of_work"`
	} `yaml:"accounts"`

	Composer struct {
		Editor string `yaml:"editor"`
	} `yaml:"composer"`

	Server struct {
		AllowInsecure bool `yaml:"allow_insecure"`
	} `yaml:"server"`

	Resolver struct {
		Local struct {
			Path string `yaml:"path"`
		} `yaml:"local"`

		Remote struct {
			URL string `yaml:"url"`
		} `yaml:"remote"`
	} `yaml:"resolver"`
}

// LoadConfig loads the client configuration from the given path
func (c *ClientConfig) LoadConfig(configPath string) error {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var lc wrappedClientConfig = wrappedClientConfig{}
	err = yaml.Unmarshal(data, &lc)
	if err != nil {
		return err
	}

	// We only care about the Cfg section. This keeps our "config:" section in the yaml file but we can still use
	// config.Client.Logger.Level instead of config.Client.Cfg.Logger.Level
	*c = lc.Cfg

	return nil
}
