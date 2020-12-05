package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

// OpenJSONFileEditor will open a text editor where you can manually edit the given src json
func OpenJSONFileEditor(src interface{}, dst interface{}) error {
	// Create a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "bmtmpedit-")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Write our json data to the temp file
	data, err := json.MarshalIndent(src, "", "  ")
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
