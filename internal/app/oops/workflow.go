package oops

import (
	"fmt"
	"strings"

	"github.com/silahicamil/lgtm/internal/app/util"
)

type CommitEntry struct {
	Hash    string
	Subject string
}

type OopsResult struct {
	CurrentBranch string
	Commits       []CommitEntry
}

func (o *OopsResult) GetCurrentBranch() error {
	branch, err := util.RunGit("branch", "--show-current")
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	o.CurrentBranch = branch
	return nil
}

// GetRecentCommits fetches the last n commits as short hash + subject lines.
func (o *OopsResult) GetRecentCommits(n int) error {
	output, err := util.RunGit("log", fmt.Sprintf("-%d", n), "--pretty=format:%h %s")
	if err != nil {
		return fmt.Errorf("failed to get git log: %w", err)
	}

	o.Commits = nil
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		subject := ""
		if len(parts) == 2 {
			subject = parts[1]
		}
		o.Commits = append(o.Commits, CommitEntry{
			Hash:    parts[0],
			Subject: subject,
		})
	}
	return nil
}

// BuildDisplayList returns formatted strings for each commit for the prompt UI.
func (o *OopsResult) BuildDisplayList() []string {
	items := make([]string, len(o.Commits))
	for i, c := range o.Commits {
		items[i] = fmt.Sprintf("%d. %s - %s", i+1, c.Hash, c.Subject)
	}
	return items
}

// ResetToCommit performs a soft reset to the commit before the selected one,
// keeping changes in the working directory.
func (o *OopsResult) ResetToCommit(idx int) error {
	if idx < 0 || idx >= len(o.Commits) {
		return fmt.Errorf("invalid commit index: %d", idx)
	}
	hash := o.Commits[idx].Hash
	_, err := util.RunGit("reset", "--soft", hash)
	if err != nil {
		return fmt.Errorf("failed to reset to %s: %w", hash, err)
	}
	return nil
}

// ForcePush force-pushes the current branch to origin after a reset.
func (o *OopsResult) ForcePush() error {
	_, err := util.RunGit("push", "origin", o.CurrentBranch, "--force-with-lease")
	if err != nil {
		return fmt.Errorf("failed to force push: %w", err)
	}
	return nil
}
