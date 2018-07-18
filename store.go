package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessfraz/pony/gpg"
)

// SecretFile is the structure for how the
// decrypted secret filestorage is organized.
type SecretFile struct {
	Secrets map[string]string `json:"secrets,omitempty"`
}

// readSecretsFile opens the secrets filestore,
// decrypts the file contents, and unmarshals
// the contents as SecretsFile.
func readSecretsFile(filename string) (s SecretFile, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return s, fmt.Errorf("opening secrets file failed: %v", err.Error())
	}
	defer f.Close()

	// decrypt the file
	decryptedFile, err := gpg.Decrypt(f, filepath.Join(gpgpath, "secring.gpg"), keyid)
	if err != nil {
		return s, fmt.Errorf("gpg decrypt file failed: %v", err)
	}

	// unmarshal the contents
	jsonParser := json.NewDecoder(decryptedFile)
	if err = jsonParser.Decode(&s); err != nil {
		return s, fmt.Errorf("json decoding decrypted file failed: %v", err)
	}

	return s, err
}

// writeSecretsFile takes a SecretsFile struct marshals,
// encrypts, and writes it to the secrets filestore.
func writeSecretsFile(filename string, s SecretFile) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Could not create filestore for secrets at %s: %v", file, err)
	}
	defer f.Close()

	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling secret file to json failed: %v", s)
	}

	// encrypt the string to the file
	encryptedBytes, err := gpg.Encrypt(b, filepath.Join(gpgpath, "pubring.gpg"))
	if err != nil {
		return fmt.Errorf("gpg encrypt on file create failed: %v", err)
	}

	if _, err := f.Write(encryptedBytes); err != nil {

		return fmt.Errorf("writing to file failed: %v", err)
	}

	return nil
}

func (s *SecretFile) setKeyValue(key, value string, force bool) error {
	// add the key value pair to secrets
	if len(s.Secrets) == 0 {
		s.Secrets = map[string]string{
			key: value,
		}
	} else {
		// check if the key already exists
		// and warn the user we are overwriting
		if val, ok := s.Secrets[key]; ok && !force {
			return fmt.Errorf("Secret for (%s) already exists with value (%s), use `update` command instead", key, val)
		}
		s.Secrets[key] = value
	}

	return writeSecretsFile(file, *s)
}
