package hop

import (
	"fmt"
	"strings"

	"github.com/silahicamil/lgtm/internal/app/util"
)

type HopResult struct {
	CurrentBranch string
	Branches      []string
}

func (h *HopResult) GetBranches() error {
	output, err := util.RunGit("branch")
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	h.Branches = nil
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "* ") {
			h.CurrentBranch = strings.TrimPrefix(line, "* ")
		} else {
			h.Branches = append(h.Branches, line)
		}
	}
	return nil
}

func (h *HopResult) Checkout(branch string) error {
	_, err := util.RunGit("checkout", branch)
	if err != nil {
		return fmt.Errorf("failed to checkout %s: %w", branch, err)
	}
	return nil
}
