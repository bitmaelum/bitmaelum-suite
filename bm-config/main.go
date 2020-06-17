package main

import (
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/bm-config/cmd"
    "os"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Path to your configuration file"`
    Version     bool        `short:"v" long:"version" description:"Display version information"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    if opts.Version {
       core.WriteVersionInfo("BitMaelum Server", os.Stdout)
       fmt.Println()
       os.Exit(1)
    }

    fmt.Println(core.GetAsciiLogo())

    core.LoadClientConfigOrPass(opts.Config)
    core.LoadServerConfigOrPass(opts.Config)

    cmd.Execute()
}
