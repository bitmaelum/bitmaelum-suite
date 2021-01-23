package cmd

import (
	"fmt"
	"os"
)

func fatal(args ...interface{}) {
	fmt.Println(append([]interface{}{"Error: "}, args...)...)
	os.Exit(1)
}

func warn(args ...interface{}) {
	fmt.Println(append([]interface{}{"Warning: "}, args...)...)
}
