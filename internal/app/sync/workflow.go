package sync

import (
	"fmt"

	"github.com/silahicamil/lgtm/internal/app/util"
)

type SyncResult struct {
	CurrentBranch string
	OriginBranch  string
	Stash         bool // Whether to stash or not
	StashPos      int  // Stash position (0 usually)
}

func (sr *SyncResult) GetCurrentBranch() error {
	branch, err := util.RunGit("branch", "--show-current")
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}
	sr.CurrentBranch = branch
	return nil
}

func (sr *SyncResult) FetchOrigin() error {
	_, err := util.RunGit("fetch", "origin", sr.OriginBranch)
	if err != nil {
		return fmt.Errorf("failed to fetch origin/%s: %w", sr.OriginBranch, err)
	}
	return nil
}

func (sr *SyncResult) Merge() error {
	remote := "origin/" + sr.OriginBranch
	_, err := util.RunGit("merge", remote)
	if err != nil {
		return fmt.Errorf("failed to merge %s into %s: %w", remote, sr.CurrentBranch, err)
	}
	return nil
}

func (sr *SyncResult) StashChanges() error {
	_, err := util.RunGit("stash")
	if err != nil {
		return fmt.Errorf("failed to stash changes: %w", err)
	}
	sr.Stash = true
	return nil
}

func (sr *SyncResult) StashPop() error {
	_, err := util.RunGit("stash", "pop")
	if err != nil {
		return fmt.Errorf("failed to pop stash: %w", err)
	}
	sr.Stash = false
	return nil
}
