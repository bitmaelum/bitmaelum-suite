// Copyright (c) 2022 BitMaelum Authors
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

package organisation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	errIncorrectKeyFormat    = errors.New("incorrect key format")
	errIncorrectKeyType      = errors.New("incorrect key typ")
	errUnknownValidationType = errors.New("unknown validation type")
)

const (
	// TypeDNS is the validation through DNS TXT records
	TypeDNS = "dns"
	// TypeKeyBase is the validation through KeyBase
	TypeKeyBase = "kb"
	// TypeGPG is the validation through GPG/PGP
	TypeGPG = "gpg"
)

// ValidationType defines a validation method and data for validating an organisation in different ways
type ValidationType struct {
	Type  string
	Value string
}

// NewValidationTypeFromStringArray generates validation types based on an array of strings
func NewValidationTypeFromStringArray(arr []string) ([]ValidationType, error) {
	vals := []ValidationType{}

	for _, s := range arr {
		v, err := NewValidationTypeFromString(s)
		if err != nil {
			return nil, errors.New(s)
		}

		vals = append(vals, *v)
	}

	return vals, nil
}

// NewValidationTypeFromString creates a new validation type from the given string
func NewValidationTypeFromString(s string) (*ValidationType, error) {
	v := &ValidationType{}

	if !strings.Contains(s, " ") {
		return nil, errIncorrectKeyFormat
	}

	// <type> <data>
	parts := strings.SplitN(s, " ", 2)

	// Check type
	switch strings.ToLower(parts[0]) {
	case TypeDNS:
		v.Type = TypeDNS
	case TypeGPG:
		v.Type = TypeGPG
	case TypeKeyBase:
		v.Type = TypeKeyBase
	default:
		return nil, errIncorrectKeyType
	}

	// Set value
	val := strings.TrimSpace(parts[1])
	if len(val) == 0 {
		return nil, errIncorrectKeyFormat
	}
	v.Value = parts[1]

	return v, nil
}

func (v *ValidationType) String() string {
	return fmt.Sprintf("%s %s", v.Type, v.Value)
}

// MarshalJSON marshals a key into bytes
func (v *ValidationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON unmarshals bytes into a key
func (v *ValidationType) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	v1, err := NewValidationTypeFromString(s)
	if err != nil {
		return err
	}

	// This seems wrong, but copy() doesn't work?
	v.Type = v1.Type
	v.Value = v1.Value

	return err
}

// Validate will validate the given validation type and returns true when its validated correctly
func (v *ValidationType) Validate(o Organisation) (bool, error) {
	switch v.Type {
	case TypeDNS:
		return validateDNS(o, v.Value)
	default:
		return false, errUnknownValidationType
	}
}
