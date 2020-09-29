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
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/password"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"path/filepath"
)

const (
	pbkdfIterations = 100002
)

const (
	// VersionV1 is the first version that uses versioning
	VersionV1 = iota + 1
)

// VaultPassword is the given password through the commandline for opening the vault
var VaultPassword string

// VaultData hold the actual data that is encrypted inside the vault
type VaultData struct { //nolint
	Accounts      []internal.AccountInfo      `json:"accounts"`
	Organisations []internal.OrganisationInfo `json:"organisations"`
}

// Vault defines our vault with path and password. Only the accounts should be exported
type Vault struct {
	Version  int // Version of the vault data. Should increase on changes so we can easily migrate if needed
	Data     VaultData
	RawData  []byte
	password []byte
	path     string
}

// vaultContainer is a json wrapper that encrypts the actual vault data
type vaultContainer struct {
	Data []byte `json:"data"`
	Salt []byte `json:"salt"`
	Iv   []byte `json:"iv"`
	Hmac []byte `json:"hmac"`
}

// New instantiates a new vault
func New(p string, pwd []byte) (*Vault, error) {
	var err error

	v := &Vault{
		Version: VersionV1,
		Data: VaultData{
			Accounts:      []internal.AccountInfo{},
			Organisations: []internal.OrganisationInfo{},
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
	_, err = os.Stat(p)
	if _, ok := err.(*os.PathError); ok {
		err = os.MkdirAll(filepath.Dir(p), 0777)
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

	// // Decrypt the encrypted vault and return its data
	// data, err = v.readEncryptedData()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// // Populate the vault based on the data
	// v, err := decodeVaultJSON(data)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// switch v.(type) {
	// case OldVault:
	// 	v.Data.Accounts = v.Accounts
	// default:
	// 	logrus.Error("unknown vault data encountered")
	// 	return nil, err
	// }
	//
	//
	//
	// // Otherwise, read vault data, if possible
	// err = v.readVault()
	// return v, err
}

// // readVault decrypts and reads the vault
// func (v *Vault) readVault() error {
// 	data, err := readFileData(v.path)
// 	if err != nil {
// 		return err
// 	}
//
// 	container := &vaultContainer{}
// 	err = json.Unmarshal(data, &container)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Check if HMAC is correct
// 	hash := hmac.New(sha256.New, v.password)
// 	hash.Write(container.Data)
// 	if bytes.Compare(hash.Sum(nil), container.Hmac) != 0 {
// 		return errors.New("incorrect password")
// 	}
//
// 	// Generate key based on password
// 	derivedAESKey := pbkdf2.Key(v.password, container.Salt, pbkdfIterations, 32, sha256.New)
// 	aes256, err := aes.NewCipher(derivedAESKey)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Decrypt vault data
// 	plainText := make([]byte, len(container.Data))
// 	ctr := cipher.NewCTR(aes256, container.Iv)
// 	ctr.XORKeyStream(plainText, container.Data)
//
// 	// store raw data. This makes editing through vault-edit easier
// 	v.RawData = plainText
//
// 	v.populateVaultDataFromRaw()
// 	return nil
// }

// sanityCheck checks if the vault contains correct data. It might be the accounts are in some kind of invalid state,
// so we should not save any data once we detected this.
func (v *Vault) sanityCheck() bool {
	if len(v.Data.Accounts) == 0 {
		return false
	}

	for _, acc := range v.Data.Accounts {
		if acc.PrivKey.S == "" {
			return false
		}
		if acc.PubKey.S == "" {
			return false
		}
	}

	for _, org := range v.Data.Organisations {
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
	if !v.sanityCheck() {
		return errors.New("Vault seems to have invalid data. Refusing to overwrite the current vault")
	}

	container, err := v.EncryptContainer()
	if err != nil {
		return err
	}

	// Make backup of the vault for now
	err = os.Rename(v.path, v.path+".backup")
	if err != nil {
		return err
	}

	// Write vault container back through temp file
	return internal.WriteFileWithLock(v.path, container, 0600)
}

// ReadFromDisk will read the account data from disk and stores this into the vault data
func (v *Vault) ReadFromDisk() error {
	data, err := internal.ReadFileWithLock(v.path)
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

// DecryptContainer decrypts a container and fills the values in v.Data
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
	return nil
}

// EncryptContainer encrypts v.Data and returns the vault as encrypted JSON container
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
	plainText, err := json.MarshalIndent(&v.Data, "", "  ")
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
		Data: cipherText,
		Salt: salt,
		Iv:   iv,
		Hmac: hash.Sum(nil),
	}, "", "  ")
}

// GetAccountOrDefault find the address from the vault. If address is empty, it will fetch the default address, or the
// first address in the vault if no default address is present.
func GetAccountOrDefault(vault *Vault, a string) *internal.AccountInfo {
	if a == "" {
		return vault.GetDefaultAccount()
	}

	addr, err := address.New(a)
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}

	info, err := vault.GetAccountInfo(*addr)
	if err != nil {
		logrus.Fatal("Address not found in vault")
		os.Exit(1)
	}

	return info
}

// OpenVault returns an opened vault, or opens the vault, asking a password if needed
func OpenVault() *Vault {
	fromVault := false
	if VaultPassword == "" {
		VaultPassword, fromVault = password.AskPassword()
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
		_ = password.StorePassword(VaultPassword)
	}

	return v
}
