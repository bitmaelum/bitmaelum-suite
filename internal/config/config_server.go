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

package config

import (
	"io"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
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
		Format  string `yaml:"log_format" default:"text"`
		LogPath string `yaml:"log_path" default:"stdout"`

		ApacheLogging bool   `yaml:"apache_log" default:"false"`
		ApacheLogPath string `yaml:"apache_log_path"`
	} `yaml:"logging"`

	Work struct {
		Pow struct {
			Bits int `yaml:"bits"`
		} `yaml:"pow"`
	} `yaml:"work"`

	Paths struct {
		Processing string `yaml:"processing"`
		Retry      string `yaml:"retry"`
		Incoming   string `yaml:"incoming"`
		Accounts   string `yaml:"accounts"`
	} `yaml:"paths"`

	Server struct {
		Hostname      string `yaml:"hostname"`
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		CertFile      string `yaml:"certfile"`
		KeyFile       string `yaml:"keyfile"`
		VerboseInfo   bool   `yaml:"verbose_info"`
		AllowInsecure bool   `yaml:"allow_insecure"`
		RoutingFile   string `yaml:"routingfile"`
	} `yaml:"server"`

	Management struct {
		Enabled bool `yaml:"remote_enabled"`
	} `yaml:"management"`

	Organisations []string `yaml:"organisations"`

	Webhooks struct {
		Enabled bool   `yaml:"enabled"`
		System  string `yaml:"system"`
		Workers int    `yaml:"workers"`
	} `yaml:"webhooks"`

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

	Bolt struct {
		DatabasePath string `yaml:"database_path"`
	} `yaml:"bolt"`

	DefaultResolver string `yaml:"default_resolver"`
	Resolvers       struct {
		Remote struct {
			URL           string `yaml:"url"`
			AllowInsecure bool   `yaml:"allow_insecure"`
		} `yaml:"remote"`
		Sqlite struct {
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
		Chain []string `yaml:"chain"`
	} `yaml:"resolvers"`
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
	c.Server.CertFile, _ = homedir.Expand(c.Server.CertFile)
	c.Server.KeyFile, _ = homedir.Expand(c.Server.KeyFile)
	c.Server.RoutingFile, _ = homedir.Expand(c.Server.RoutingFile)
	c.Bolt.DatabasePath, _ = homedir.Expand(c.Bolt.DatabasePath)
	c.Resolvers.Sqlite.Path, _ = homedir.Expand(c.Resolvers.Sqlite.Path)

	loadedConfig = "server"

	return nil
}
