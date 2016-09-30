package gpg

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/term"
	"golang.org/x/crypto/openpgp"
)

func readInput(in io.Reader, out io.Writer) []byte {
	reader := bufio.NewReader(in)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintln(out, err.Error())
		os.Exit(1)
	}
	return line
}

// Decrypt a io.Reader with the given secretKeyring.
// You can optionally pass a defaultGPGKey to use for the
// decryption, otherwise it will use the first entity.
func Decrypt(f io.Reader, secretKeyring, defaultGPGKey string) (io.Reader, error) {
	// Open the private key file
	keyringFileBuffer, err := os.Open(secretKeyring)
	if err != nil {
		return nil, err
	}
	defer keyringFileBuffer.Close()

	entityList, err := openpgp.ReadKeyRing(keyringFileBuffer)
	if err != nil {
		return nil, err
	}

	var entity *openpgp.Entity
	if defaultGPGKey != "" {

		// loop through their keys until we find the one they want
		var foundKey bool
		for _, e := range entityList {
			// we can match on the fingerprint or the keyid because
			// why not? I bet no one knows the difference
			if e.PrimaryKey.KeyIdString() == defaultGPGKey ||
				e.PrimaryKey.KeyIdShortString() == defaultGPGKey ||
				fmt.Sprintf("%X", e.PrimaryKey.Fingerprint) == defaultGPGKey {
				foundKey = true
				entity = e
				break
			}
		}

		if !foundKey {
			// we didn't find the key they specified
			return nil, fmt.Errorf("Could not find private GPG Key with id: %s", defaultGPGKey)
		}

	} else {
		// they didn't set a default key
		// so let's hope it is the first one :/
		// TODO(jessfraz): maybe prompt here if they have
		// more than one private key
		entity = entityList[0]
	}

	var identityString string
	for _, identity := range entity.Identities {
		identityString = fmt.Sprintf(" %s [%s]", identity.Name, entity.PrimaryKey.KeyIdString())
		break
	}

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	stdin, _, stderr := term.StdStreams()
	stdinFd, _ := term.GetFdInfo(stdin)
	oldState, err := term.SaveState(stdinFd)
	if err != nil {
		return nil, err
	}

	// prompt for passphrase
	fmt.Fprintf(stderr, "GPG Passphrase for key%s: ", identityString)
	term.DisableEcho(stdinFd, oldState)

	// read what they inputed
	passphrase := readInput(stdin, stderr)
	fmt.Fprint(stderr, "\n\n")

	// restore the terminal
	term.RestoreTerminal(stdinFd, oldState)

	logrus.Debugln("Decrypting private key using passphrase")

	entity.PrivateKey.Decrypt(passphrase)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphrase)
	}
	logrus.Debugln("Finished decrypting private key using passphrase")

	// base64 decode
	dec := base64.NewDecoder(base64.StdEncoding, f)

	// Decrypt it with the contents of the private key
	md, err := openpgp.ReadMessage(dec, entityList, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("GPG ReadMessage failed: %v", err)
	}

	// return the body contents
	return md.UnverifiedBody, nil
}
