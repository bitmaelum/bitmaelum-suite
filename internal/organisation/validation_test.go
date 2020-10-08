package organisation

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
)

func TestNewValidationTypeFromString(t *testing.T) {
	var (
		v   *ValidationType
		err error
	)

	v, err = NewValidationTypeFromString("")
	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = NewValidationTypeFromString("dns")
	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = NewValidationTypeFromString("dns         ")
	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = NewValidationTypeFromString("unknown foobar.org")
	assert.Error(t, err)
	assert.Nil(t, v)

	v, err = NewValidationTypeFromString("dns foobar.org")
	assert.NoError(t, err)
	assert.Equal(t, TypeDNS, v.Type)
	assert.Equal(t, "foobar.org", v.Value)

	v, err = NewValidationTypeFromString("kb bitmaelum")
	assert.NoError(t, err)
	assert.Equal(t, TypeKeyBase, v.Type)
	assert.Equal(t, "bitmaelum", v.Value)

	v, err = NewValidationTypeFromString("gpg 0xDEADBEEF")
	assert.NoError(t, err)
	assert.Equal(t, TypeGPG, v.Type)
	assert.Equal(t, "0xDEADBEEF", v.Value)
}

func TestValidationType_MarshalJSON(t *testing.T) {
	v, _ := NewValidationTypeFromString("dns foobar.org")

	j, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.Equal(t, "\"dns foobar.org\"", string(j))

	v1 := &ValidationType{}
	err = json.Unmarshal(j, &v1)
	assert.NoError(t, err)
	assert.Equal(t, "dns", v1.Type)
	assert.Equal(t, "foobar.org", v1.Value)

	err = json.Unmarshal([]byte("\"asdfasfdsafasf\""), &v1)
	assert.Error(t, err)

	err = json.Unmarshal([]byte("\"\""), &v1)
	assert.Error(t, err)

	err = json.Unmarshal([]byte("\"unknown foobar.org\""), &v1)
	assert.Error(t, err)
}

func TestValidationType_String(t *testing.T) {
	v, _ := NewValidationTypeFromString("dns foobar.org")
	assert.Equal(t, "dns foobar.org", v.String())

	v.Value = "bitmaelum.com"
	v.Type = TypeGPG
	assert.Equal(t, "gpg bitmaelum.com", v.String())
}

func TestValidate(t *testing.T) {
	v, _ := NewValidationTypeFromString("dns bitmaelum.org")

	o := &Organisation{
		Hash:       hash.New("bitmaelum"),
		FullName:   "BitMaelum",
		PublicKey:  bmcrypto.PubKey{},
		Validation: []ValidationType{*v},
	}

	// No error, no result
	resolver = &mockResolver{}
	ok, err := v.Validate(*o)
	assert.NoError(t, err)
	assert.False(t, ok)

	// No error, correct result
	resolver.SetCallbackTXT(func() ([]string, error) {
		return []string{
			"00000000004a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
			"49aa67181f4a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
		}, nil
	})
	ok, err = v.Validate(*o)
	assert.NoError(t, err)
	assert.True(t, ok)

	// No error, incorrect result
	resolver.SetCallbackTXT(func() ([]string, error) {
		return []string{
			"00000000004a3176f9b65605390bb81126e8ff1f6d03b1bd220c53e7a6b36d3e",
		}, nil
	})

	ok, err = v.Validate(*o)
	assert.NoError(t, err)
	assert.False(t, ok)

	// Error, no result
	resolver.SetCallbackTXT(func() ([]string, error) {
		return []string{
			"",
		}, errors.New("foobar")
	})

	ok, err = v.Validate(*o)
	assert.Error(t, err)
	assert.False(t, ok)
}
