package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\U0001F4E8 Enter address: ")
		addr, _ := reader.ReadString('\n')
		addr = strings.Trim(addr, "\n")

		if addr == "" {
			break
		}

		a, err := address.NewAddress(addr)
		if err != nil {
			panic(err)
		}
		fmt.Printf("NEW: %s\n", a.Hash().String())
	}
}
