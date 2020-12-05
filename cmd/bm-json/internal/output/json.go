package output

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JSONT map[string]interface{}

// JSONOut outputs a specific interface
func JSONOut(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}

	fmt.Print(string(b))
}


// JSONErrorOut outputs an error
func JSONErrorOut(err error) {
	v := map[string]interface{}{
		"error": err.Error(),
	}

	JSONOut(v)
}


// JSONErrorStrOut outputs an error string
func JSONErrorStrOut(s string) {
	JSONErrorOut(errors.New(s))
}
