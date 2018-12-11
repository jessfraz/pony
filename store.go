package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessfraz/pony/gpg"
)

// secretFile is the structure for how the decrypted secret filestorage is organized.
type secretFile struct {
	Secrets map[string]string `json:"secrets,omitempty"`
}

// readSecretsFile opens the secrets filestore, decrypts the file contents,
// and unmarshals the contents as SecretsFile.
func readSecretsFile(filename string) (s secretFile, err error) {
	f, err := filepath.Abs(filename)
	if err != nil {
		return s, err
	}

	// Decrypt the file.
	file, err := gpg.Decrypt(f)
	if err != nil {
		return s, fmt.Errorf("gpg decrypt file failed: %v", err)
	}

	// Unmarshal the contents.
	if err = json.NewDecoder(file).Decode(&s); err != nil {
		return s, err
	}

	return s, err
}

// writeSecretsFile takes a SecretsFile struct and marshals, encrypts, and
// writes it to the secrets filestore.
func writeSecretsFile(filename string, s secretFile) error {
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling secret file to json failed: %v", s)
	}

	// Encrypt the string to the file
	eb, err := gpg.Encrypt(b, keyid)
	if err != nil {
		return fmt.Errorf("gpg encrypt on file create failed: %v", err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not create filestore for secrets at %s: %v", file, err)
	}
	defer f.Close()

	if _, err := f.Write(eb); err != nil {
		return fmt.Errorf("writing to file failed: %v", err)
	}

	return nil
}

func (s *secretFile) setKeyValue(key, value string, force bool) error {
	// Add the key value pair to secrets.
	if len(s.Secrets) == 0 {
		s.Secrets = map[string]string{
			key: value,
		}
	} else {
		// Check if the key already exists and warn the user we are overwriting.
		if _, ok := s.Secrets[key]; ok && !force {
			return fmt.Errorf("Secret for key %s already exists, use `--force` to overwrite", key)
		}
		s.Secrets[key] = value
	}

	return writeSecretsFile(file, *s)
}
