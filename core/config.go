package core

import (
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/mitchellh/go-homedir"
)

// Load client configuration from given path
func LoadClientConfig(path string) {
    p, _ := homedir.Expand(path)
    err := config.Client.LoadConfig(p)
    if err != nil {
        panic(err)
    }
}

// Load server configuration from given path
func LoadServerConfig(path string) {
    p, _ := homedir.Expand(path)
    err := config.Server.LoadConfig(p)
    if err != nil {
        panic(err)
    }
}
