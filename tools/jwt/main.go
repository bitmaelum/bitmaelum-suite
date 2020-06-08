package main

import (
    "fmt"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/account"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
    Addr        string      `short:"a" long:"address" description:"account"`
    Password    string      `short:"p" long:"password" description:"Password to decrypt your account"`
}

var opts Options


func main() {
    core.ParseOptions(&opts)
    core.LoadClientConfig(opts.Config)

    // Convert strings into addresses
    fromAddr, err := core.NewAddressFromString(opts.Addr)
    if err != nil {
        panic(err)
    }

    // Load account
    var pwd = []byte(opts.Password)
    ai, err := account.LoadAccount(*fromAddr, pwd)
    if err != nil {
       panic(err)
    }

    ts,err := core.GenerateJWTToken(fromAddr.Hash(), ai.PrivKey)
    if err != nil {
        panic(err)
    }

    fmt.Println(ts)

    v, err := core.ValidateJWTToken(ts, fromAddr.Hash(), ai.PubKey)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%#v\n", v)
}
