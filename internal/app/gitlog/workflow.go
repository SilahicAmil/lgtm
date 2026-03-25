package gitlog

import (
	"fmt"
	"strings"

	"github.com/silahicamil/lgtm/internal/app/util"
)

type LogResult struct {
	Lines []string
}

func (l *LogResult) Fetch(limit int) error {
	output, err := util.RunGit(
		"log",
		fmt.Sprintf("-%d", limit),
		"--oneline",
		"--graph",
		"--decorate",
	)
	if err != nil {
		return fmt.Errorf("failed to get git log: %w", err)
	}

	l.Lines = nil
	for _, line := range strings.Split(output, "\n") {
		if line != "" {
			l.Lines = append(l.Lines, line)
		}
	}
	return nil
}
