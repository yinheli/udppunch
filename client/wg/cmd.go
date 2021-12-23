package wg

import (
	"os/exec"
	"strings"
)

func run(cmd string, args ...string) (string, error) {
	b, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
