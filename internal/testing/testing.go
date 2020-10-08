package testing

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// ReadTestKey reads a path to a keypair and returns the keys
func ReadTestKey(p string) (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, nil, err
	}

	type jsonKeyType struct {
		PrivKey bmcrypto.PrivKey `json:"private_key"`
		PubKey  bmcrypto.PubKey  `json:"public_key"`
	}

	v := &jsonKeyType{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, nil, err
	}

	return &v.PrivKey, &v.PubKey, nil
}

// ReadTestFile reads a file
func ReadTestFile(p string) []byte {
	data, _ := ioutil.ReadFile(p)
	return data
}

// ReadJSON reads a json file and returns it in the given interface
func ReadJSON(p string, v interface{}) error {
	data, _ := ioutil.ReadFile(p)

	return json.Unmarshal(data, v)
}
