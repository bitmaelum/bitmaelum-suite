package config

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

// Server keeps all server configuration settings
var Server = ServerConfig{}

// Basically, our config is inside the "config" section. So we load the whole file and only store the Cfg section
type wrappedServerConfig struct {
	Cfg ServerConfig `yaml:"config"`
}

// ServerConfig is the representation of the server configuration
type ServerConfig struct {
	Logging struct {
		Level   string `yaml:"log_level"`
		LogPath string `yaml:"log_path" default:"stdout"`

		ApacheLogging bool   `yaml:"apache_log" default:"false"`
		ApacheLogPath string `yaml:"apache_log_path"`
	} `yaml:"logging"`

	Accounts struct {
		ProofOfWork int `yaml:"proof_of_work"`
	} `yaml:"accounts"`

	Paths struct {
		Processing string `yaml:"processing"`
		Retry      string `yaml:"retry"`
		Incoming   string `yaml:"incoming"`
		Accounts   string `yaml:"accounts"`
	} `yaml:"paths"`

	Server struct {
		Name          string `yaml:"hostname"`
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		CertFile      string `yaml:"certfile"`
		KeyFile       string `yaml:"keyfile"`
		VerboseInfo   bool   `yaml:"verbose_info"`
		AllowInsecure bool   `yaml:"allow_insecure"`
	} `yaml:"server"`

	Management struct {
		Enabled bool `yaml:"remote_enabled"`
	} `yaml:"management"`

	Acme struct {
		Enabled         bool   `yaml:"enabled"`
		Domain          string `yaml:"domain"`
		Path            string `yaml:"path"`
		Email           string `yaml:"email"`
		RenewBeforeDays string `yaml:"renew_days"`
	} `yaml:"acme"`

	Redis struct {
		Host string `yaml:"host"`
		Db   int    `yaml:"port"`
	} `yaml:"redis"`

	Resolver struct {
		Local struct {
			Path string `yaml:"path"`
		} `yaml:"local"`

		Remote struct {
			URL string `yaml:"url"`
		} `yaml:"remote"`
	} `yaml:"resolver"`
}

// LoadConfig loads the server configuration from the given path
func (c *ServerConfig) LoadConfig(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var lc wrappedServerConfig = wrappedServerConfig{}
	err = yaml.Unmarshal(data, &lc)
	if err != nil {
		return err
	}

	// We only care about the Cfg section. This keeps our "config:" section in the yaml file but we can still use
	// config.Server.Logger.Level instead of config.Server.Cfg.Logger.Level
	*c = lc.Cfg

	// Expand homedirs in configuration
	c.Logging.LogPath, _ = homedir.Expand(c.Logging.LogPath)
	c.Logging.ApacheLogPath, _ = homedir.Expand(c.Logging.ApacheLogPath)
	c.Paths.Processing, _ = homedir.Expand(c.Paths.Processing)
	c.Paths.Retry, _ = homedir.Expand(c.Paths.Retry)
	c.Paths.Incoming, _ = homedir.Expand(c.Paths.Incoming)
	c.Paths.Accounts, _ = homedir.Expand(c.Paths.Accounts)
	c.Acme.Path, _ = homedir.Expand(c.Acme.Path)
	c.Resolver.Local.Path, _ = homedir.Expand(c.Resolver.Local.Path)
	c.Server.CertFile, _ = homedir.Expand(c.Server.CertFile)
	c.Server.KeyFile, _ = homedir.Expand(c.Server.KeyFile)

	return nil
}
