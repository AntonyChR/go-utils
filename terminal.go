package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
)

// GetTerminalSize retrieves the current terminal size (number of rows and columns).
//
// This function executes the `stty size` command to obtain the current terminal
// dimensions, where the expected output is two numbers: the number of rows (height)
// and the number of columns (width).
//
// Note: This function is only available on Linux systems. It relies on the 
// `stty` command, which may not be available or function the same way on other 
// operating systems.
//
// Returns:
//   - height (int): The number of rows in the terminal.
//   - width (int): The number of columns in the terminal.
//   - err (error): An error if any issue occurs while executing the command
//     or processing the output. It returns errors in cases such as:
//     - The command execution fails.
//     - The command output is unexpectedly empty or not in the expected format.
//     - The dimensions from the output cannot be converted to integers.
//
// Example usage:
//
//   height, width, err := GetTerminalSize()
//   if err != nil {
//       log.Fatalf("Error retrieving terminal size: %v", err)
//   }
//   fmt.Printf("Height: %d, Width: %d\n", height, width)
func GetTerminalSize() (height, width int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
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
