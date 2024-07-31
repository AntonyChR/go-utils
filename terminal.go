package utils

import (
	"bytes"
	"os"
	"os/exec"
)

// Get terminal dimentions
func GetTermDims() (height, width string, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	cmd.Stdin = os.Stdout
	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	lines := bytes.Split(out, []byte("\n"))
	splt := bytes.Split(lines[0], []byte(" "))
	return string(splt[0]), string(splt[1]), err
}
