package utils

import (
	"bytes"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func UnmarshalStrict(in []byte, out interface{}) (err error) {
	knownFieldsDecoder := yaml.NewDecoder(bytes.NewReader(in))
	knownFieldsDecoder.KnownFields(true)
	return knownFieldsDecoder.Decode(out)
}

func CopyFile(sourceFile, destinationFolder string) error {
	// Read the source file into memory.
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	// Get the destination file path.
	destinationFile := filepath.Join(destinationFolder, filepath.Base(sourceFile))

	// Check if the destination file already exists.
	_, err = os.Stat(destinationFile)
	if err == nil {
		// If the file already exists, remove it.
		err = os.Remove(destinationFile)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		// If there was an error other than the file not existing, return it.
		return err
	}

	// Write the source file to the destination folder.
	err = os.WriteFile(destinationFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
