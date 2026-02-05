package ship

import (
	"fmt"
	"os/exec"
	"strings"
)

type ShipResult struct {
	Branchname string
	CleanFiles map[string]string
	DirtyFiles map[string]string
	Completed  map[string]bool
}

// Eventually move to /git folder
// Split each into it's own reusable thing maybe?
// Probably have
func (sr *ShipResult) CheckStatusAndBranch() error {
	// Get the current branch
	branchCmd := exec.Command("git", "branch", "--show-current")
	branchOutput, err := branchCmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}
	sr.Branchname = strings.TrimSpace(string(branchOutput))

	// get modified/untracked files
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOutput, err := statusCmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	lines := strings.SplitSeq(string(statusOutput), "\n")

	// parse git stauts output
	// and add all to CleanFiles for now
	// We could add one it to like PendingFiles? instead
	// I think CleanFIles makes it just easier
	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fileName := strings.TrimSpace(line[3:])
		fmt.Println(fileName)

		sr.CleanFiles[fileName] = "Initial Scan"
	}

	return nil
}

func (sr *ShipResult) CheckDiff() (*ShipResult, error) {

	diffNameCmd := exec.Command("git", "diff", "--name-only")
	diffNameOutput, err := diffNameCmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// grep didn’t match anything — this is normal
			fmt.Println("No matching patterns found")
		} else {
			return sr, fmt.Errorf("git diff failed: %w", err)
		}
	}

	// Need the relative path or else it won't work
	// repoRootCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	// rootBytes, err := repoRootCmd.CombinedOutput()
	// repoRoot := strings.TrimSpace(string(rootBytes))

	fmt.Printf("diffOuput %s", string(diffNameOutput))
	files := strings.Split(strings.TrimSpace(string(diffNameOutput)), "\n")
	patterns := []string{"console.log", "error_log"}

	for _, file := range files {
		fmt.Printf("FIle %s", file)
		cmd := exec.Command("git", "diff", file)
		diffBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("git diff failed for file:", file, err)
			continue
		}

		currentFile := file
		for _, line := range strings.Split(string(diffBytes), "\n") {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				for _, pattern := range patterns {
					if strings.Contains(line, pattern) {
						sr.DirtyFiles[currentFile] = "Contains: " + pattern
						delete(sr.CleanFiles, currentFile)
						sr.Completed[currentFile] = true
					}
				}
			}
		}
	}

	// Run the git diff grep
	// FOr any files that get returned
	// remove it from cleanfiles
	// add to dirty files
	// then return acordingly in the CLI
	// Just for now loop over and return
	// File Name - Reason to keep it simple
	return sr, nil
}
