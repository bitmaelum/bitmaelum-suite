package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"io/ioutil"
)

type options struct {
	Password string `short:"p" long:"password" description:"Password to your vault" required:"true"`
	SrcFile  string `short:"s" long:"source" description:"JSON source for the vault" required:"true"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)

	// Create new (empty) vault
	v, err := vault.New("", []byte(opts.Password))
	if err != nil {
		panic(err)
	}

	// Read data from source
	data, err := ioutil.ReadFile(opts.SrcFile)
	if err != nil {
		panic(err)
	}
	// Unmarshal into account
	err = json.Unmarshal(data, &v.Accounts)
	if err != nil {
		panic(err)
	}

	// output encrypted vault
	buf, err := v.Encrypted()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf)
}
