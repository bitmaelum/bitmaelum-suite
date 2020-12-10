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

package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

// OpenJSONFileEditor will open a text editor where you can manually edit the given src json
func JSONFileEditor(src interface{}, dst interface{}) error {
	// Create a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "bmtmpedit-")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Write our json data to the temp file
	data, _ := json.MarshalIndent(src, "", "  ")
	_, err = tmpFile.Write(data)
	_ = tmpFile.Sync()
	if err != nil {
		return err
	}

	editor, err := getEditor()
	if err != nil {
		return err
	}

	for {
		// Execute editor
		c := exec.Command(editor, tmpFile.Name())
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		err = c.Run()
		if err != nil {
			return err
		}

		// Editor errored
		if !c.ProcessState.Success() {
			return err
		}

		data, err = ioutil.ReadFile(tmpFile.Name())
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, dst)
		if err != nil || dst == nil {
			// Something is not right with unmarshalling. Let the user try and edit the file again
			continue
		}

		// All is good, return
		return nil
	}
}

// Get default editor
func getEditor() (string, error) {
	if os.Getenv("EDITOR") != "" {
		return os.Getenv("EDITOR"), nil
	}

	editors := []string{"/usr/bin/editor", "/usr/bin/nano"}
	for _, editor := range editors {
		_, err := os.Stat(editor)
		if err == nil {
			return editor, nil
		}
	}

	return "", errors.New("no editor found")
}
