package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"io/ioutil"
	"os"
	"os/exec"
)

type options struct {
	Config      string `short:"c" long:"config" description:"Path to your configuration file"`
	Password    string `short:"p" long:"password" description:"Password to your vault"`
	NewPassword string `short:"n" long:"new-password" description:"New password to your vault"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	v, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		panic(err)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "vd-")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	_, err = tmpFile.Write(v.RawData)
	_ = tmpFile.Sync()
	if err != nil {
		panic(err)
	}

	editor := "/usr/bin/nano"
	if os.Getenv("EDITOR") != "" {
		editor = os.Getenv("EDITOR")
	}

	c := exec.Command(editor, tmpFile.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &v.Data)
	if err != nil {
		panic(err)
	}

	if opts.NewPassword != "" {
		v.ChangePassword(opts.NewPassword)
	}

	err = v.WriteToDisk()
	if err != nil {
		panic(err)
	}

	fmt.Println("Vault saved to disk")
}
