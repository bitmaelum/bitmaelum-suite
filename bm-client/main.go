package main

import (
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/bm-client/account"
    "github.com/bitmaelum/bitmaelum-server/bm-client/cmd"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/bitmaelum/bitmaelum-server/core/password"
    "os"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Path to your configuration file"`
    Password    string      `short:"p" long:"password" description:"Vault password" default:""`
    Version     bool        `short:"v" long:"version" description:"Display version information"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    if opts.Version {
       core.WriteVersionInfo("BitMaelum Client", os.Stdout)
       fmt.Println()
       os.Exit(1)
    }

    fmt.Println(core.GetAsciiLogo())
    core.LoadClientConfig(opts.Config)

    if opts.Password == "" {
        opts.Password = string(password.AskPassword())
    }

    // Unlock vault
    err := account.UnlockVault(config.Client.Accounts.Path, []byte(opts.Password))
    if err != nil {
        fmt.Printf("Error while opening vault: %s", err)
        os.Exit(1)
    }

    cmd.Execute()
}
