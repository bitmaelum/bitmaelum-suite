package fileio

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path"
)

// SaveCertFiles saves the given cert and key PEM strings into the configured cert and key file. Old files are backed
// as .001 (or .002 etc) if the files already exists.
func SaveCertFiles(certPem string, keyPem string) error {
	suffix := findHighestSuffix(config.Server.Server.CertFile, config.Server.Server.KeyFile)

	var (
		newPath string
		oldPath string
		err     error
	)

	newPath, _ = homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.CertFile, suffix))
	oldPath, _ = homedir.Expand(config.Server.Server.CertFile)
	_, err = os.Stat(oldPath)
	if err == nil {
		fmt.Printf("   - moving old cert file to %s: ", newPath)
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Println("ok")
	}

	newPath, _ = homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.KeyFile, suffix))
	oldPath, _ = homedir.Expand(config.Server.Server.KeyFile)
	_, err = os.Stat(oldPath)
	if err == nil {
		fmt.Printf("   - moving old key file to %s: ", newPath)
		err = os.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Println("ok")
	}

	fmt.Printf("   - Writing new cert file %s: ", config.Server.Server.CertFile)
	newPath, _ = homedir.Expand(config.Server.Server.CertFile)
	err = ioutil.WriteFile(newPath, []byte(certPem), 0600)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	fmt.Printf("   - Writing new key file %s: ", config.Server.Server.CertFile)
	newPath, _ = homedir.Expand(config.Server.Server.KeyFile)
	err = ioutil.WriteFile(newPath, []byte(keyPem), 0600)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	return nil
}

// FindHighestSuffix returns the highest suffix found on the files. It checks first .001, then .002 etc until it finds
// a suffix that doesn't exist on all files.
func findHighestSuffix(files ...string) int {
	var suffix = 1

	for {
		var found = false
		for _, file := range files {
			p, _ := homedir.Expand(fmt.Sprintf("%s.%03d", file, suffix))
			_, err1 := os.Stat(p)
			if err1 == nil {
				found = true
				break
			}
		}

		if !found {
			return suffix
		}

		suffix++
	}
}

// LoadFile loads and unmarshals a given file
func LoadFile(p string, v interface{}) error {
	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// SaveFile saves a structured as marshalled JSON
func SaveFile(p string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	p, err = homedir.Expand(p)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(p), 755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, data, 0600)
}
