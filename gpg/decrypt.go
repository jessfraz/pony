package gpg

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

// Decrypt a file.
func Decrypt(f string) (io.Reader, error) {
	// Open and read the file.
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	// Base64 decode.
	out, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	// Run the command.
	args := []string{"--decrypt"}
	cmd := exec.Command("gpg", args...)
	cmd.Stdin = bytes.NewBuffer(out)
	out, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gpg [gpg %s] failed with stdout %q, error: %v", strings.Join(args, " "), string(out), err)
	}
	return bytes.NewReader(out), nil
}
