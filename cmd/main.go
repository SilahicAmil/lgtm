package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/silahicamil/lgtm/internal/app/cli"
	"github.com/silahicamil/lgtm/internal/app/ship"
)

var RED_CLI_PROMPT = color.New(color.FgHiRed).Add(color.Underline)
var WHITE_UNDERLINE = color.New(color.FgWhite).Add(color.BgBlack)

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
		// fmt.Println(shipRes)

		diffd, err := shipRes.CheckDiff()
		if err != nil {
			// TODO: Make this better
			fmt.Println("here2")
		}
		// fmt.Println(diffd)

		//  TODO: Make this prettier and actually useable for the end user
		for file, match := range diffd.DirtyFiles {
			RED_CLI_PROMPT.Printf("File %s", file)
			WHITE_UNDERLINE.Printf(" - %s \n", match)
			// Ask the user if they want to continue or exit
		}

		// show ALL files
		// then put a checkbox next to each box you want to add
		// Each one of those files we want to add it to like a CommitResult struct
		// Once added we just run the git add for each file
		// or git do a git add . / * if they select ALL
		// or just exit so

		// AddFiles
		fmt.Println(cmd)
	case "sync":
		fmt.Println("Sync")
	}

}
