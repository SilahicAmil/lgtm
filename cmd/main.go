package main

import (
	"fmt"
	"os"

	"github.com/silahicamil/lgtm/internal/app/cli"
	"github.com/silahicamil/lgtm/internal/app/ship"
)

func main() {
	// Parse the CLI Flag
	// cmd, args := cli.Parse(os.Args[1:])
	// Check if there is an error
	arg := os.Args[1:]

	cmd, _, err := cli.Parse(arg)

	if err != nil {
		return
	}

	// Summary
	// Phase 1: Check we are in a git repo
	// What branch are we on?
	// Is the repo dirty?
	//
	// Phase 2: Inspect diffs/files
	// Just show what was found
	//
	// PHASE 1 and 2 DONE
	// Phase 3: Show clean vs dirty
	// Ask user what they want to do
	// Can exit early
	//
	// Phase 4: Mutate
	// git add, git commit, git push
	switch cmd {
	case "ship":
		// Check status
		// Check Diff
		// Just print the DirtyFiles for now
		shipRes := &ship.ShipResult{
			Branchname: "N/A",
			CleanFiles: make(map[string]string),
			DirtyFiles: make(map[string]string),
			Completed:  make(map[string]bool),
		}
		err := shipRes.CheckStatusAndBranch()
		if err != nil {
			fmt.Println("Here?")
		}
		fmt.Println(shipRes)

		diffd, err := shipRes.CheckDiff()
		if err != nil {
			fmt.Println("here2")
		}
		fmt.Println(diffd)

		for file, match := range diffd.DirtyFiles {
			fmt.Printf("File %s - Contains: %s \n", file, match)
		}
		fmt.Println(cmd)
	case "sync":
		fmt.Println("Sync")
	}

}
