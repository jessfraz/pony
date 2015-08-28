package gpg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
)

// Encrypt a byte with the given publicKeyring
func Encrypt(b []byte, publicKeyring string) ([]byte, error) {
	// Read in public key
	keyringFileBuffer, err := os.Open(publicKeyring)
	if err != nil {
		return nil, fmt.Errorf("opening public key %s failed: %v", publicKeyring, err)
	}
	defer keyringFileBuffer.Close()

	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	// encrypt string
	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, entityList, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	// write the byte to our encrypt writer
	if _, err = w.Write(b); err != nil {
		return nil, err
	}

	// close the writer
	if err = w.Close(); err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}

	// base64 encode
	encStr := base64.StdEncoding.EncodeToString(bytes)
	return []byte(encStr), nil
}
