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

package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/console"
	"github.com/spf13/afero"
)

var (
	errIncorrectPassword   = errors.New("incorrect password")
	errNotOverwritingVault = errors.New("vault seems to have invalid data. Refusing to overwrite the current vault")
	errVaultNotFound       = errors.New("vault not found")
	errNoPathSet           = errors.New("vault path is not set")
)

// Vault versions
const (
	VersionV0 = iota // Only accounts
	VersionV1        // Accounts + organisations
	VersionV2        // Multi key
)

// LatestVaultVersion Everything below this version will be automatically migrated to this version
const LatestVaultVersion = VersionV2

// Override for testing purposes
var fs = afero.NewOsFs()

var (
	// VaultPassword is the given password through the commandline for opening the vault
	VaultPassword string
	// VaultPath is the default vault path
	VaultPath string
)

// StoreType hold the actual data that is encrypted inside the vault
type StoreType struct {
	Accounts      []AccountInfo      `json:"accounts"`
	Organisations []OrganisationInfo `json:"organisations"`
}

// StoreTypeV1 hold the actual data that is encrypted inside the vault
type StoreTypeV1 struct {
	Accounts      []AccountInfoV1      `json:"accounts"`
	Organisations []OrganisationInfoV1 `json:"organisations"`
}

// Vault defines our vault with path and password. Only the accounts should be exported
type Vault struct {
	Store    StoreType
	RawData  []byte
	password string
	path     string
}

// New instantiates a new vault
func New() *Vault {
	return &Vault{
		Store: StoreType{
			Accounts:      []AccountInfo{},
			Organisations: []OrganisationInfo{},
		},
		RawData: []byte{},
	}
}

// NewPersistent instantiates a new vault and persists on disk
func NewPersistent(p, pass string) *Vault {
	v := New()
	v.SetPassword(pass)
	v.SetPath(p)

	return v
}

// sanityCheck checks if the vault contains correct data. It might be the accounts are in some kind of invalid state,
// so we should not save any data once we detected this.
func (v *Vault) sanityCheck() bool {
	for _, acc := range v.Store.Accounts {
		for _, k := range acc.Keys {
			if k.PrivKey.S == "" {
				return false
			}
			if k.PubKey.S == "" {
				return false
			}
		}
	}

	for _, org := range v.Store.Organisations {
		for _, k := range org.Keys {
			if k.PrivKey.S == "" {
				return false
			}
			if k.PubKey.S == "" {
				return false
			}
		}
	}

	return true
}

// Persist saves the vault data back to disk
func (v *Vault) Persist() error {
	if v.path == "" {
		return errNoPathSet
	}

	// Only do sanity check when file is already present
	_, err := fs.Stat(v.path)
	fileExists := err == nil

	if fileExists && !v.sanityCheck() {
		return errNotOverwritingVault
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
	return afero.WriteFile(fs, v.path, container, 0600)
}

// SetPassword allows us to change the vault password. Will take effect on writing to disk
func (v *Vault) SetPassword(pass string) {
	v.password = pass
}

// SetPath sets the path of the vault.
func (v *Vault) SetPath(p string) {
	v.path = p
}

// Create will create a new vault on the given path
func Create(p, pass string) (*Vault, error) {
	v := NewPersistent(p, pass)

	err := v.Persist()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// OpenDefaultVault returns an opened vault on vault.VaultPath and with password vault.VaultPath. Will die when incorrect vault or password
func OpenDefaultVault() *Vault {
	if !Exists(VaultPath) {
		fmt.Printf("Cannot open vault at '%s'. You might want to initialize a new vault with by issuing \n\n   $ bm-client vault init\n", VaultPath)
		fmt.Println("")
		os.Exit(1)
	}

	var retrievedFromKeyStore = false

	// Ask password when none is found
	if VaultPassword == "" {
		VaultPassword, retrievedFromKeyStore = console.AskPassword()
	}

	// If the password was correct and not already read from the vault, store it in the vault
	if !retrievedFromKeyStore {
		_ = console.StorePassword(VaultPassword)
	}

	return OpenOrDie(VaultPath, VaultPassword)
}

// OpenOrDie will open a specific vault with a specific password
func OpenOrDie(vp, pass string) *Vault {
	v, err := Open(vp, pass)
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	return v
}

// Open will open a specific vault with a specific password
func Open(vp, pass string) (*Vault, error) {
	if !Exists(vp) {
		return nil, errVaultNotFound
	}

	// Create in memory vault
	v := NewPersistent(vp, pass)

	// Read data container from path
	data, err := afero.ReadFile(fs, v.path)
	if err != nil {
		return nil, err
	}

	container := &EncryptedContainer{}
	err = json.Unmarshal(data, &container)
	if err != nil {
		return nil, err
	}

	// Decrypt the container into the new vault
	err = v.DecryptContainer(container)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Exists will return true if the vault exists
func Exists(p string) bool {
	info, err := fs.Stat(p)

	if err != nil {
		return false
	}

	return !info.IsDir()
}
