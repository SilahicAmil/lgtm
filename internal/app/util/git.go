package util

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunGit executes a git command with the given arguments and returns
func RunGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w\n%s", args[0], err, output)
	}
	return strings.TrimSpace(string(output)), nil
}
