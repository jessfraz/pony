package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jessfraz/pony/gpg"
)

// SecretFile is the structure for how the
// decrypted secret filestorage is organized.
type SecretFile struct {
	Secrets map[string]string `json:"secrets,omitempty"`
}

// preChecks makes sure the user has gpg set up for saving secrets
// as well as a filestore file created (even if blank).
// That way we wont have to be sure about making sure the file exists later.
func preChecks() error {
	gpgErrorString := "Have you generated a gpg key? You can do so with `$ gpg --gen-key`."

	if _, err := os.Stat(publicKeyring); os.IsNotExist(err) {

		return fmt.Errorf("GPG Public Keyring (%s) does not exist.\n%s", publicKeyring, gpgErrorString)
	}

	if _, err := os.Stat(secretKeyring); os.IsNotExist(err) {

		return fmt.Errorf("GPG Secret Keyring (%s) does not exist.\n%s", secretKeyring, gpgErrorString)
	}

	// create our secrets filestore if it does not exist
	if _, err := os.Stat(filestore); os.IsNotExist(err) {
		if err := writeSecretsFile(filestore, SecretFile{}); err != nil {
			return err
		}
	}

	return nil
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
	decryptedFile, err := gpg.Decrypt(f, secretKeyring, defaultGPGKey)
	if err != nil {
		return s, fmt.Errorf("gpg decrypt file failed: %v", err)
	}

	// unmarshal the contents
	jsonParser := json.NewDecoder(decryptedFile)
	if err = jsonParser.Decode(&s); err != nil {
		return s, fmt.Errorf("json decoding decrypted file failed: %v", err.Error())
	}

	return s, err
}

// writeSecretsFile takes a SecretsFile struct marshals,
// encrypts, and writes it to the secrets filestore.
func writeSecretsFile(filename string, s SecretFile) error {
	f, err := os.Create(filestore)
	if err != nil {
		return fmt.Errorf("Could not create filestore for secrets at %s: %v", filestore, err)
	}
	defer f.Close()

	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshaling secret file to json failed: %v", s)
	}

	// encrypt the string to the file
	encryptedBytes, err := gpg.Encrypt(b, publicKeyring)
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

	return writeSecretsFile(filestore, *s)
}
