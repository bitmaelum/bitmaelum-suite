package account

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/juju/fslock"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const (
	pbkdfIterations = 100002
)

// Vault holds all our account information
var Vault *VaultType // UnlockVault initializes this value

// VaultType represents our vault
type VaultType struct {
	Accounts []core.AccountInfo
	password []byte
	path     string
}

type vaultJSONData struct {
	Data []byte `json:"data"`
	Salt []byte `json:"salt"`
	Iv   []byte `json:"iv"`
	Hmac []byte `json:"hmac"`
}

// UnlockVault takes a path and password to initialize the vault. It will try and unlock the vault and store the
// information into Vault.Accounts
func UnlockVault(p string, pwd []byte) error {
	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	Vault = &VaultType{
		Accounts: []core.AccountInfo{},
		password: pwd,
		path:     p,
	}
	err = Vault.unlockVault()
	if _, ok := err.(*os.PathError); ok {
		err = os.MkdirAll(path.Dir(p), 0777)
		if err != nil {
			return err
		}
		err = Vault.Save()
	}

	if err != nil {
		return err
	}

	return nil
}

// unlockVault unlocks the vault by the given password
func (v *VaultType) unlockVault() error {
	// Lock vault for reading
	lock := fslock.New(v.path + ".lock")
	err := lock.TryLock()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(v.path)
	_ = lock.Unlock()
	_ = os.Remove(v.path + ".lock")
	if err != nil {
		return err
	}

	vaultData := &vaultJSONData{}
	err = json.Unmarshal(data, &vaultData)
	if err != nil {
		return err
	}

	// Check if HMAC is correct
	hash := hmac.New(sha256.New, v.password)
	hash.Write(vaultData.Data)
	if bytes.Compare(hash.Sum(nil), vaultData.Hmac) != 0 {
		return errors.New("incorrect password")
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key(v.password, vaultData.Salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return err
	}

	// Decrypt vault data
	plainText := make([]byte, len(vaultData.Data))
	ctr := cipher.NewCTR(aes256, vaultData.Iv)
	ctr.XORKeyStream(plainText, vaultData.Data)

	// Unmarshal vault data
	var accounts []core.AccountInfo
	err = json.Unmarshal(plainText, &accounts)
	if err != nil {
		return err
	}

	v.Accounts = accounts
	return nil
}

// Lock Locks the vault by removing the password and account data
func (v *VaultType) Lock() {
	v.password = nil
	v.Accounts = nil
}

// Add adds a new account to the vault
func (v *VaultType) Add(account core.AccountInfo) {
	v.Accounts = append(v.Accounts, account)
}

// Remove the given account from the vault
func (v *VaultType) Remove(address core.Address) {
	k := 0
	for _, acc := range v.Accounts {
		if acc.Address != address.String() {
			v.Accounts[k] = acc
			k++
		}
	}
	v.Accounts = v.Accounts[:k]
}

// Save saves the account data back into the vault on disk
func (v *VaultType) Save() error {
	if v.password == nil {
		return errors.New("vault is locked and cannot be saved")
	}

	// Generate 64 byte salt
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key(v.password, salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return err
	}

	// Generate 32 byte IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	// Marshal and encrypt the data
	plainText, err := json.MarshalIndent(&v.Accounts, "", "  ")
	if err != nil {
		return err
	}

	cipherText := make([]byte, len(plainText))
	ctr := cipher.NewCTR(aes256, iv)
	ctr.XORKeyStream(cipherText, plainText)

	// Generate HMAC based on the encrypted data (encrypt-then-mac?)
	hash := hmac.New(sha256.New, v.password)
	hash.Write(cipherText)

	// Generate the vault structure for disk
	data, err := json.MarshalIndent(&vaultJSONData{
		Data: cipherText,
		Salt: salt,
		Iv:   iv,
		Hmac: hash.Sum(nil),
	}, "", "  ")
	if err != nil {
		return err
	}

	// Write vault back through temp file
	return safeWrite(v.path, data, 0600)
}

// FindAccount tries to find the given address and returns the account from the vault
func (v *VaultType) FindAccount(address core.Address) (*core.AccountInfo, error) {
	for _, acc := range v.Accounts {
		if acc.Address == address.String() {
			return &acc, nil
		}
	}

	return nil, errors.New("cannot find account")
}

// HasAccount returns true when the vault has an account for the given address
func (v *VaultType) HasAccount(address core.Address) bool {
	_, err := v.FindAccount(address)

	return err == nil
}

// safeWrite writes data by safely writing to a temp file first
func safeWrite(path string, data []byte, perm os.FileMode) error {
	// Lock the file first. Make sure we are the only one working on it
	lock := fslock.New(path + ".lock")
	err := lock.TryLock()
	if err != nil {
		return err
	}

	defer func() {
		_ = lock.Unlock()
		_ = os.Remove(path + ".lock")
	}()

	err = ioutil.WriteFile(path+".tmp", data, perm)
	if err != nil {
		return err
	}

	err = os.Rename(path+".tmp", path)
	return err
}
