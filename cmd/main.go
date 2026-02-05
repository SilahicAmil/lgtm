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
