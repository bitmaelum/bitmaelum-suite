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
	"io"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Client keeps all client configuration settings
var Client ClientConfig = ClientConfig{}

// Basically, our config is inside the "config" section. So we load the whole file and only store the Cfg section
type wrappedClientConfig struct {
	Cfg ClientConfig `yaml:"config"`
}

// ClientConfig is the representation of the client configuration
type ClientConfig struct {
	Vault struct {
		Path string `yaml:"path"`
	} `yaml:"vault"`

	Composer struct {
		Editor string `yaml:"editor"`
	} `yaml:"composer"`

	Server struct {
		AllowInsecure bool `yaml:"allow_insecure"`
		DebugHTTP     bool `yaml:"debug_http"`
	} `yaml:"server"`

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

// LoadConfig loads the client configuration from the given path
func (c *ClientConfig) LoadConfig(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
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

	// Expand homedirs in configuration
	c.Vault.Path, _ = homedir.Expand(c.Vault.Path)
	c.Resolvers.Sqlite.Path, _ = homedir.Expand(c.Resolvers.Sqlite.Path)

	loadedConfig = "client"

	return nil
}
