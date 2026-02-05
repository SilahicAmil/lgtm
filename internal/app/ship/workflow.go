package ship

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ShipResult struct {
	Branchname string
	CleanFiles map[string]string
	DirtyFiles map[string]string
	Completed  map[string]bool
}

var patterns = []string{"console.log", "error_log"}

// TODO: Eventually move to /git folder
// Split each into it's own reusable thing maybe?
// Probably have like a client.go to handle the exec.command stuff
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

		fileName := strings.TrimSpace(line[2:])
		sr.CleanFiles[fileName] = "All Good"
	}

	return nil
}

func (sr *ShipResult) CheckDiff() (*ShipResult, error) {

	diffNameCmd := exec.Command("git", "diff", "--name-only")
	diffNameOutput, err := diffNameCmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// grep didn’t match anything so just keep going
			fmt.Println("No matching patterns found")
		} else {
			return sr, fmt.Errorf("git diff failed: %w", err)
		}
	}

	// Need the relative path or else it won't work
	// So just go to working directory
	repoRootCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	rootBytes, err := repoRootCmd.Output()
	repoRoot := strings.TrimSpace(string(rootBytes))
	os.Chdir(repoRoot)

	files := strings.Split(strings.TrimSpace(string(diffNameOutput)), "\n")

	if err != nil {
		fmt.Println("Loading Patterns File Failed: %w", err)
	}

	// Loop over the files from the --name-only
	for _, file := range files {
		// Run a diff on each file
		cmd := exec.Command("git", "diff", file)
		diffBytes, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("git diff failed for file:", file, err)
			continue
		}

		currentFile := file
		// Verify it is an actual diff chunk
		for _, line := range strings.Split(string(diffBytes), "\n") {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				for _, pattern := range patterns {
					if strings.Contains(line, string(pattern)) {
						// Add to Dirty Files
						// Remove from CleanFiles
						sr.DirtyFiles[currentFile] = "Contains: " + string(pattern)
						delete(sr.CleanFiles, currentFile)
						// sr.Completed[currentFile] = true
					}
				}
			}
		}
	}
	return sr, nil
}
