// Copyright (c) 2020 BitMaelum Authors
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

package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/spf13/afero"
	"golang.org/x/crypto/pbkdf2"
)

const (
	pbkdfIterations = 100002
)

// VersionV0 is the first version that uses versioning
const (
	VersionV0 = iota
	VersionV1
)

// Override for testing purposes
var fs = afero.NewOsFs()

// VaultPassword is the given password through the commandline for opening the vault
var VaultPassword string

// StoreType hold the actual data that is encrypted inside the vault
type StoreType struct {
	Accounts      []AccountInfo      `json:"accounts"`
	Organisations []OrganisationInfo `json:"organisations"`
}

// Vault defines our vault with path and password. Only the accounts should be exported
type Vault struct {
	Store    StoreType
	RawData  []byte
	password []byte
	path     string
}

// vaultContainer is a json wrapper that encrypts the actual vault data
type vaultContainer struct {
	Version int    `json:"version"`
	Data    []byte `json:"data"`
	Salt    []byte `json:"salt"`
	Iv      []byte `json:"iv"`
	Hmac    []byte `json:"hmac"`
}

// New instantiates a new vault
func New(p string, pwd []byte) (*Vault, error) {
	var err error

	v := &Vault{
		Store: StoreType{
			Accounts:      []AccountInfo{},
			Organisations: []OrganisationInfo{},
		},
		RawData:  []byte{},
		password: pwd,
		path:     p,
	}

	// No path given, we return just the empty vault
	if p == "" {
		return v, nil
	}

	// Create new vault when we cannot find the one specified
	_, err = fs.Stat(p)
	if _, ok := err.(*os.PathError); ok {
		err = fs.MkdirAll(filepath.Dir(p), 0777)
		if err != nil {
			return nil, err
		}
		err = v.WriteToDisk()
		return v, err
	}

	// Read vault data from disk
	err = v.ReadFromDisk()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// sanityCheck checks if the vault contains correct data. It might be the accounts are in some kind of invalid state,
// so we should not save any data once we detected this.
func (v *Vault) sanityCheck() bool {
	if len(v.Store.Accounts) == 0 {
		return false
	}

	for _, acc := range v.Store.Accounts {
		if acc.PrivKey.S == "" {
			return false
		}
		if acc.PubKey.S == "" {
			return false
		}
	}

	for _, org := range v.Store.Organisations {
		if org.PrivKey.S == "" {
			return false
		}
		if org.PubKey.S == "" {
			return false
		}
	}

	return true
}

// WriteToDisk saves the vault data back to disk
func (v *Vault) WriteToDisk() error {
	// Only do sanity chck when file is already present
	_, err := fs.Stat(v.path)
	fileExists := err == nil

	if fileExists && !v.sanityCheck() {
		return errors.New("vault seems to have invalid data. Refusing to overwrite the current vault")
	}

	container, err := v.EncryptContainer()
	if err != nil {
		return err
	}

	// Make backup of the vault for now
	if fileExists {
		err = fs.Rename(v.path, v.path+".backup")
		if err != nil {
			return err
		}
	}

	// Write vault container back
	err = afero.WriteFile(fs, v.path, container, 0600)
	return err
}

// ReadFromDisk will read the account data from disk and stores this into the vault data
func (v *Vault) ReadFromDisk() error {
	data, err := afero.ReadFile(fs, v.path)
	if err != nil {
		return err
	}

	container := &vaultContainer{}
	err = json.Unmarshal(data, &container)
	if err != nil {
		return err
	}

	return v.DecryptContainer(container)
}

// DecryptContainer decrypts a container and fills the values in v.Store
func (v *Vault) DecryptContainer(container *vaultContainer) error {

	// Check if HMAC is correct
	hash := hmac.New(sha256.New, v.password)
	hash.Write(container.Data)
	if !bytes.Equal(hash.Sum(nil), container.Hmac) {
		return errors.New("incorrect password")
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key(v.password, container.Salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return err
	}

	// Decrypt vault data
	plainText := make([]byte, len(container.Data))
	ctr := cipher.NewCTR(aes256, container.Iv)
	ctr.XORKeyStream(plainText, container.Data)

	// store raw data. This makes editing through vault-edit tool easier
	v.RawData = plainText

	if container.Version == VersionV0 {
		// Unmarshal "old" style, with no organisations present
		var accounts []AccountInfo
		err = json.Unmarshal(plainText, &accounts)
		if err == nil {
			v.Store.Accounts = accounts
			v.Store.Organisations = []OrganisationInfo{}

			// Write back to disk in a newer format
			return v.WriteToDisk()
		}
	}

	// Version 1 has organisation info
	if container.Version == VersionV1 {
		err = json.Unmarshal(plainText, &v.Store)
		if err != nil {
			return err
		}
	}

	return nil
}

// EncryptContainer encrypts v.Store and returns the vault as encrypted JSON container
func (v *Vault) EncryptContainer() ([]byte, error) {
	// Generate 64 byte salt
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key(v.password, salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return nil, err
	}

	// Generate 32 byte IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	// Marshal and encrypt the data
	plainText, err := json.MarshalIndent(&v.Store, "", "  ")
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))
	ctr := cipher.NewCTR(aes256, iv)
	ctr.XORKeyStream(cipherText, plainText)

	// Generate HMAC based on the encrypted data (encrypt-then-mac?)
	hash := hmac.New(sha256.New, v.password)
	hash.Write(cipherText)

	// Generate the vault structure for disk
	return json.MarshalIndent(&vaultContainer{
		Version: VersionV1,
		Data:    cipherText,
		Salt:    salt,
		Iv:      iv,
		Hmac:    hash.Sum(nil),
	}, "", "  ")
}

// ChangePassword allows us to change the vault password. Will take effect on writing to disk
func (v *Vault) ChangePassword(newPassword string) {
	v.password = []byte(newPassword)
}

// FindShortRoutingID will find a short routing ID in the vault and expand it to the full routing ID. So we can use
// "12345" instead of "1234567890123456789012345678901234567890".
// Will not return anything when multiple candidates are found.
func (v *Vault) FindShortRoutingID(id string) string {
	var found = ""
	for _, acc := range v.Store.Accounts {
		if strings.HasPrefix(acc.RoutingID, id) {
			// Found something else that matches
			if found != "" && found != acc.RoutingID {
				// Multiple entries are found, don't return them
				return ""
			}
			found = acc.RoutingID
		}
	}

	return found
}

// GetAccountOrDefault find the address from the vault. If address is empty, it will fetch the default address, or the
// first address in the vault if no default address is present.
func GetAccountOrDefault(vault *Vault, a string) *AccountInfo {
	if a == "" {
		return vault.GetDefaultAccount()
	}

	acc, _ := GetAccount(vault, a)
	return acc
}

// GetAccount returns the given account, or nil when not found
func GetAccount(vault *Vault, a string) (*AccountInfo, error) {
	addr, err := address.NewAddress(a)
	if err != nil {
		return nil, err
	}

	return vault.GetAccountInfo(*addr)
}

// OpenVault returns an opened vault, or opens the vault, asking a password if needed
func OpenVault() *Vault {
	fromVault := false

	// Check if vault exists
	if vaultExists(config.Client.Accounts.Path) {
		if VaultPassword == "" {
			VaultPassword, fromVault = console.AskPassword()
		}
	} else {
		if VaultPassword == "" {
			fmt.Printf("A vault could not be found. Creating a new vault at '%s'.\n", config.Client.Accounts.Path)
			b, err := console.AskDoublePassword()
			VaultPassword = string(b)
			if err != nil {
				fmt.Printf("Error while creating vault: %s", err)
				fmt.Println("")
				os.Exit(1)
			}
		}
	}

	// Unlock vault
	v, err := New(config.Client.Accounts.Path, []byte(VaultPassword))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	// If the password was correct and not already read from the vault, store it in the vault
	if !fromVault {
		_ = console.StorePassword(VaultPassword)
	}

	return v
}

func vaultExists(p string) bool {
	_, err := os.Stat(config.Client.Accounts.Path)
	return err == nil
}
