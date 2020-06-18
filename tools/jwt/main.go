package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/account/server"
	"log"
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
        log.Fatal(err)
    }

    // Load account
    var pwd = []byte(opts.Password)
    ai, err := server.LoadAccount(*fromAddr, pwd)
    if err != nil {
       log.Fatal(err)
    }

    fmt.Println("PRIV:")
    fmt.Printf("%s\n", ai.PrivKey)

    fmt.Println("PUB:")
    fmt.Printf("%s\n", ai.PubKey)
    fmt.Println("")

    ts,err := core.GenerateJWTToken(fromAddr.Hash(), ai.PrivKey)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(ts)

    token, err := core.ValidateJWTToken(ts, fromAddr.Hash(), ai.PubKey)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%#v\n", token.Claims)
}
