package main

import (
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/bm-config/cmd"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Path to your configuration file"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    fmt.Println(core.GetAsciiLogo())

    core.LoadClientConfigOrPass(opts.Config)
    core.LoadServerConfigOrPass(opts.Config)

    cmd.Execute()
}
