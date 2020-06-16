package main

import (
    "github.com/bitmaelum/bitmaelum-server/core"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
    Addr        string      `short:"a" long:"address" description:"account"`
    Password    string      `short:"p" long:"password" description:"Password to decrypt your account"`
    Box         string      `short:"b" long:"box" description:"message box"`
    Id          string      `short:"i" long:"id" description:"message id"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    core.LoadClientConfig(opts.Config)


    //// Convert strings into addresses
    //fromAddr, err := core.NewAddressFromString(opts.Addr)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //
    //var pwd = []byte(opts.Password)

    //// Load our FROM account
    //ai, err := account.LoadAccount(*fromAddr, pwd)
    //if err != nil {
    //    log.Fatal(err)
    //}

    //as := container.GetAccountService()
    //mbi := as.FetchMessageBoxes(fromAddr.Hash(), opts.Box)
    //if err != nil {
    //    log.Fatal(err)
    //}


    //message, err := message.LoadMessage(fromAddr, opts.Id)
    //if err != nil {
    //    log.Fatal(err)
    //}

    // Do things with message
}
