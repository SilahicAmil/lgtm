package ship

import (
	"fmt"
	"os"
	"strings"

	"github.com/silahicamil/lgtm/internal/app/util"
)

// TODO: I need to make this more generic
// and actually seperate each thing

type ShipResult struct {
	Branchname string
	CleanFiles map[string]string
	DirtyFiles map[string]string
	Completed  map[string]bool
}

type CommitSelection struct {
	Branchname    string
	SelectedFiles map[string]string // filename -> status/reason
	IncludesDirty bool              // true if user chose to include flagged files
	CommitMessage string
}

// NewCommitSelection creates an initialized CommitSelection
func NewCommitSelection() *CommitSelection {
	return &CommitSelection{
		Branchname:    "",
		SelectedFiles: make(map[string]string),
		IncludesDirty: false,
		CommitMessage: "",
	}
}

// SelectAll merges both CleanFiles and DirtyFiles into SelectedFiles
func (cs *CommitSelection) SelectAll(sr *ShipResult) {
	for file, status := range sr.CleanFiles {
		cs.SelectedFiles[file] = status
	}
	for file, status := range sr.DirtyFiles {
		cs.SelectedFiles[file] = status
		cs.IncludesDirty = true
	}
}

// SelectFiles adds specific files to the selection from ShipResult
func (cs *CommitSelection) SelectFiles(sr *ShipResult, files []string) {
	for _, file := range files {
		if status, ok := sr.CleanFiles[file]; ok {
			cs.SelectedFiles[file] = status
		} else if status, ok := sr.DirtyFiles[file]; ok {
			cs.SelectedFiles[file] = status
			cs.IncludesDirty = true
		}
	}
}

// GetAllFilesList returns a combined list of all files (clean + dirty) for selection UI
func (sr *ShipResult) GetAllFilesList() []string {
	files := make([]string, 0, len(sr.CleanFiles)+len(sr.DirtyFiles))
	for file := range sr.CleanFiles {
		files = append(files, file)
	}
	for file := range sr.DirtyFiles {
		files = append(files, file)
	}
	return files
}

// IsDirtyFile checks if a file is in the DirtyFiles map
func (sr *ShipResult) IsDirtyFile(file string) bool {
	_, ok := sr.DirtyFiles[file]
	return ok
}

func readPatternsFile() ([]byte, error) {
	return os.ReadFile("config/patterns.txt")
}

func (sr *ShipResult) CheckBranch() error {
	// Get the current branch
	branchOutput, err := util.RunGit("branch", "--show-current")
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}
	sr.Branchname = branchOutput

	return nil
}

// TODO: Eventually move to /git folder
// Split each into it's own reusable thing maybe?
// Probably have like a client.go to handle the exec.command stuff
func (sr *ShipResult) CheckStatus() error {

	// get modified/untracked files
	statusOutput, err := util.RunGit("status", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	lines := strings.SplitSeq(statusOutput, "\n")

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
	var patterns []string

	diffNameOutput, err := util.RunGit("diff", "--name-only")
	if err != nil {
		fmt.Println("No matching patterns found")
	}

	// Need the relative path or else it won't work
	// So just go to working directory
	repoRoot, err := util.RunGit("rev-parse", "--show-toplevel")
	if err != nil {
		return sr, fmt.Errorf("failed to get repo root: %w", err)
	}
	os.Chdir(repoRoot)

	files := strings.Split(diffNameOutput, "\n")
	data, err := readPatternsFile()

	if err != nil {
		return sr, fmt.Errorf("unable to load patterns.txt file: %w", err)
	}

	str := strings.TrimSpace(string(data))
	patterns = strings.Split(str, "\n")

	// Loop over the files from the --name-only
	for _, file := range files {
		// Run a diff on each file
		diffOutput, err := util.RunGit("diff", file)
		if err != nil {
			fmt.Println("git diff failed for file:", file, err)
			continue
		}

		currentFile := file
		// Verify it is an actual diff chunk
		for _, line := range strings.Split(diffOutput, "\n") {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				for _, pattern := range patterns {
					if line == "" {
						continue
					}
					if strings.Contains(line, pattern) {
						// Add to Dirty Files
						// Remove from CleanFiles
						sr.DirtyFiles[currentFile] += " contains: " + string(pattern)
						delete(sr.CleanFiles, currentFile)
						// sr.Completed[currentFile] = true
					}
				}
			}
		}
	}
	return sr, nil
}

func (cs *CommitSelection) AddGitFiles() (string, error) {
	args := []string{"add"}

	for file := range cs.SelectedFiles {
		args = append(args, file)
	}

	_, err := util.RunGit(args...)
	if err != nil {
		return "", err
	}

	return "Added files successfully", nil
}

func (cs *CommitSelection) AddCommitMessage(commitMsg string) (string, error) {
	cs.CommitMessage = commitMsg

	_, err := util.RunGit("commit", "-m", commitMsg)
	if err != nil {
		return "", err
	}

	return "Added your message.", nil
}

func PushGit(branchName string) (string, error) {
	_, err := util.RunGit("push", "origin", branchName)
	if err != nil {
		// Push failed — undo the commit but keep changes staged
		util.RunGit("reset", "--soft", "HEAD~1")
		// Unstage everything
		util.RunGit("restore", "--staged", ".")
		return "", err
	}

	return "Pushed to branch. Thanks for breaking Prod!", nil
}
