package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

var Server ServerConfig = ServerConfig{}

// Basically, our config is inside the "config" section. So we load the whole file and only store the Cfg section
type WrappedServerConfig struct {
    Cfg ServerConfig `yaml:"config"`
}

type ServerConfig struct {
    Logging struct {
       Level    string `yaml:"level"`
    } `yaml:"logging"`

    Account struct {
       Registration bool    `yaml:"registration"`
       Path         string  `yaml:"path"`
       ProofOfWork  int     `yaml:"proof_of_work"`
    } `yaml:"account"`

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

    Resolve struct {
       Local struct {
           Path  string  `yaml:"path"`
       } `yaml:"local"`

       Remote struct {
           Url  string  `yaml:"url"`
       } `yaml:"remote"`
    } `yaml:"resolve"`
}

// Load server configuration
func (c *ServerConfig) LoadConfig(configPath string) error {
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return err
    }

    var lc WrappedServerConfig = WrappedServerConfig{}
    err = yaml.Unmarshal(data, &lc)
    if err != nil {
        return err
    }

    // We only care about the Cfg section. This keeps our "config:" section in the yaml file but we can still use
    // config.Server.Logger.Level instead of config.Server.Cfg.Logger.Level
    *c = lc.Cfg

    return nil
}

