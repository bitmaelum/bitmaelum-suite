package core

import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "os"
)

// Parse config options. Will exit when issues arrise
func ParseOptions(opts interface{}) {
    parser := flags.NewParser(opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        flagsError, _ := err.(*flags.Error)
        if flagsError.Type == flags.ErrHelp {
            os.Exit(1)
        }

        fmt.Println()
        parser.WriteHelp(os.Stdout)
        fmt.Println()
        os.Exit(1)
    }
}
