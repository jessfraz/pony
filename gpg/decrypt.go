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

func Decrypt(f io.Reader, secretKeyring string) (io.Reader, error) {
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
	logrus.Debugf("Entity List: %+v", entityList)
	// TODO(jfrazelle): find which key is good
	entity := entityList[1]

	var identityString string
	for _, identity := range entity.Identities {
		identityString = fmt.Sprintf(" %q", identity.Name)
		break
	}

	// Get the passphrase and read the private key.
	// Have not touched the encrypted string yet
	stdin, stdout, stderr := term.StdStreams()
	stdinFd, _ := term.GetFdInfo(stdin)
	oldState, err := term.SaveState(stdinFd)
	if err != nil {
		return nil, err
	}

	// prompt for passphrase
	fmt.Fprintf(stdout, "GPG Passphrase for key%s: ", identityString)
	term.DisableEcho(stdinFd, oldState)

	// read what they inputed
	passphrase := readInput(stdin, stderr)
	fmt.Fprintln(stdout, "\n")

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
