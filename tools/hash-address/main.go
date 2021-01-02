// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
		fmt.Printf("Hash: %s\n", a.Hash().String())
		fmt.Printf(" Local: %s\n", a.LocalHash().String())
		fmt.Printf(" Org: %s\n", a.OrgHash().String())
	}
}
