package organisation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// NewValidationTypeFromString creates a new validation type from the given string
func NewValidationTypeFromString(s string) (*ValidationType, error) {
	v := &ValidationType{}

	if !strings.Contains(s, " ") {
		return nil, errors.New("incorrect key format")
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
		return nil, errors.New("incorrect key type")
	}

	// Set value
	val := strings.TrimSpace(parts[1])
	if len(val) == 0 {
		return nil, errors.New("empty value")
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
		return false, errors.New("incorrect validation type")
	}
}
