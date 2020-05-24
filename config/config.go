package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

var Configuration Config = Config{}

// Basically, our config is inside the "config" section. So we load the whole file and only store the Cfg section
type LoadConfig struct {
    Cfg Config `yaml:"config"`
}

type Config struct {
    Logging struct {
       Level    string `yaml:"level"`
    } `yaml:"logging"`

    Registration struct {
       Enabled bool `yaml:"enabled"`
    } `yaml:"registration"`

    Server struct {
       Host    string   `yaml:"host"`
       Port    int      `yaml:"port"`
    } `yaml:"server"`

    TLS struct {
       CertFile   string `yaml:"certfile"`
       KeyFile    string `yaml:"keyfile"`
    } `yaml:"tls"`

    Redis struct {
       Host    string `yaml:"host"`
       Db      int `yaml:"port"`
    } `yaml:"redis"`
}

func (c *Config) LoadConfig(configPath string) error {
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return err
    }

    var lc LoadConfig = LoadConfig{}
    err = yaml.Unmarshal(data, &lc)
    if err != nil {
        return err
    }

    // We only care about the Cfg section. This keeps our "config:" section in the yaml file but we can still use
    // config.Configuration.Logger.Level instead of config.Configuration.Cfg.Logger.Level
    *c = lc.Cfg

    return nil
}
