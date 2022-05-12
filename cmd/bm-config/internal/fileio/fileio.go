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

package fileio

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/afero"
)

var fs = afero.NewOsFs()

// SaveCertFiles saves the given cert and key PEM strings into the configured cert and key file. Old files are backed
// as .001 (or .002 etc) if the files already exists.
func SaveCertFiles(certPem string, keyPem string) error {
	suffix := findHighestSuffix(config.Server.Server.CertFile, config.Server.Server.KeyFile)

	var (
		newPath string
		oldPath string
		err     error
	)

	newPath = fmt.Sprintf("%s.%03d", config.Server.Server.CertFile, suffix)
	oldPath = config.Server.Server.CertFile
	_, err = fs.Stat(oldPath)
	if err == nil {
		fmt.Printf("   - moving old cert file to %s: ", newPath)
		err := fs.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Println("ok")
	}

	newPath = fmt.Sprintf("%s.%03d", config.Server.Server.KeyFile, suffix)
	oldPath = config.Server.Server.KeyFile
	_, err = fs.Stat(oldPath)
	if err == nil {
		fmt.Printf("   - moving old key file to %s: ", newPath)
		err = fs.Rename(oldPath, newPath)
		if err != nil {
			return err
		}
		fmt.Println("ok")
	}

	fmt.Printf("   - Writing new cert file %s: ", config.Server.Server.CertFile)
	newPath = config.Server.Server.CertFile
	err = afero.WriteFile(fs, newPath, []byte(certPem), 0600)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	fmt.Printf("   - Writing new key file %s: ", config.Server.Server.CertFile)
	newPath = config.Server.Server.KeyFile
	err = afero.WriteFile(fs, newPath, []byte(keyPem), 0600)
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
			p := fmt.Sprintf("%s.%03d", file, suffix)
			_, err1 := fs.Stat(p)
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
	data, err := afero.ReadFile(fs, p)
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

	err = fs.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, p, data, 0600)
}
