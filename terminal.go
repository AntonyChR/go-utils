package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
)

func GetTermDims() (height, width int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	cmd.Stdin = os.Stdout
	out, err := cmd.Output()
	if err != nil {
		return 0,0, err 
	}

	if len(out) == 0 {
		return 0, 0, errors.New("empty output from stty command")
	}
	lines := bytes.Split(out, []byte("\n"))
	if len(lines) == 0 {
		return 0, 0, errors.New("unexpected output format from stty command")
	}
	
	parts := bytes.Split(lines[0], []byte(" "))

	height, err = strconv.Atoi(string(parts[0]))
	if err != nil {
		return 0,0, err 
	}

	width, err = strconv.Atoi(string(parts[1]))
	if err != nil {
		return 0,0, err 
	}

	return height, width, nil
}
