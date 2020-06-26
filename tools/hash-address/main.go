package main

import (
	"bufio"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var address string
	for {
		fmt.Print("\U0001F4E8 Enter address: ")
		address, _ = reader.ReadString('\n')
		address = strings.Trim(address, "\n")

		if address == "" {
			break
		}

		hash := core.StringToHash(address)
		fmt.Printf("%s\n", hash.String())
	}
}
