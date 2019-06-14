package cliutils

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// GetUserInput asks the user for input
func GetUserInput() (string, error) {
	r := bufio.NewReader(os.Stdin)
	v, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	// trim spaces from the input
	v = strings.TrimSpace(v)

	return v, nil
}

// GetPasswordInput returns a string that was SAFE
func GetPasswordInput() (string, error) {
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// GetYesOrNoInput wraps around GetUserInput but asks yes or no,
// defaults to false unless yes or y (true)
func GetYesOrNoInput() (bool, error) {
	v, err := GetUserInput()
	if err != nil {
		return false, err
	}

	v = strings.ToLower(v)

	if v == "yes" || v == "y" {
		return true, nil
	}

	return false, nil
}
