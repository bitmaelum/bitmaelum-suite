package core

import (
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/mitchellh/go-homedir"
)

// Load client configuration from given path or panic if cannot load
func LoadClientConfig(path string) {
    p, _ := homedir.Expand(path)
    err := config.Client.LoadConfig(p)
    if err != nil{
        panic(err)
    }
}

// Load server configuration from given path or panic if cannot load
func LoadServerConfig(path string) {
    p, _ := homedir.Expand(path)
    err := config.Server.LoadConfig(p)
    if err != nil {
        panic(err)
    }
}

// Load client configuration, but don't panic if we can't
func LoadClientConfigOrPass(path string) {
    p, _ := homedir.Expand(path)
    _ = config.Client.LoadConfig(p)
}

// Load server configuration, but don't panic if we can't
func LoadServerConfigOrPass(path string) {
    p, _ := homedir.Expand(path)
    _ = config.Server.LoadConfig(p)
}

