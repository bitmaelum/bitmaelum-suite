package main

import (
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/tools/mail-server-config/cmd"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./server-config.yml"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    //core.LoadServerConfig(opts.Config)

    cmd.Execute()
}
