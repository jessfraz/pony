package gpg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

// Encrypt a byte with the given publicKeyring.
func Encrypt(b []byte, keyid string) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	var stdout bytes.Buffer
	args := []string{"--encrypt"}
	if len(keyid) > 0 {
		args = append(args, "--recipient", keyid)
	}
	cmd := exec.Command("gpg", args...)
	cmd.Stdin = buf
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gpg [gpg %s] failed with stdout %q, error: %v", strings.Join(args, " "), stdout.String(), err)
	}

	// Base64 encode.
	return []byte(base64.StdEncoding.EncodeToString(out)), nil
}
