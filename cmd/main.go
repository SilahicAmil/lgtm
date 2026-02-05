package main

import (
	"fmt"
	"os"

	"github.com/silahicamil/lgtm/internal/app/cli"
	"github.com/silahicamil/lgtm/internal/app/ship"
)

func main() {
	// Parse the CLI Flag
	arg := os.Args[1:]

	cmd, _, err := cli.Parse(arg)

	if err != nil {
		return
	}

	switch cmd {
	case "ship":
		shipRes := &ship.ShipResult{
			Branchname: "N/A",
			CleanFiles: make(map[string]string),
			DirtyFiles: make(map[string]string),
			Completed:  make(map[string]bool),
		}
		err := shipRes.CheckStatusAndBranch()
		if err != nil {
			// TODO: Make this better
			fmt.Println("Here?")
		}
		fmt.Println(shipRes)

		diffd, err := shipRes.CheckDiff()
		if err != nil {
			// TODO: Make this better
			fmt.Println("here2")
		}
		fmt.Println(diffd)

		//  TODO: Make this prettier and actually useable for the end user
		for file, match := range diffd.DirtyFiles {
			fmt.Printf("File %s - Contains: %s \n", file, match)
		}
		fmt.Println(cmd)
	case "sync":
		fmt.Println("Sync")
	}

}
