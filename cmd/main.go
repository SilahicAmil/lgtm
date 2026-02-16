package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/silahicamil/lgtm/internal/app/cli"
	"github.com/silahicamil/lgtm/internal/app/ship"
	"github.com/spaceweasel/promptui"
)

var RED_CLI_PROMPT = color.New(color.FgHiRed).Add(color.Underline)
var GREEN_CLI_PROMPT = color.New(color.FgGreen).Add(color.BgBlack)
var WHITE_UNDERLINE = color.New(color.FgWhite).Add(color.BgBlack)

var commands = map[string]string{
	"help":  "Show this help message",
	"ship":  "Stage, commit, and push files",
	"sync":  "Sync your branch with another one",
	"oops":  "Undo last commit and reset changes",
	"quote": "Inspirational quote to get you through the day",
}

func main() {
	// Parse the CLI Flag
	arg := os.Args[1:]

	cmd, _, err := cli.Parse(arg)

	if err != nil {
		return
	}

	// Start a reader

	switch cmd {
	case "ship":
		shipRes := &ship.ShipResult{
			Branchname: "N/A",
			CleanFiles: make(map[string]string),
			DirtyFiles: make(map[string]string),
			Completed:  make(map[string]bool),
		}
		err := shipRes.CheckBranch()
		if err != nil {
			fmt.Println("checkbranch err", err)
		}

		err = shipRes.CheckBranch()
		if err != nil {
			// TODO: Make this better
			fmt.Println("Here?", err)
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
			RED_CLI_PROMPT.Printf("File - %s", file)
			WHITE_UNDERLINE.Printf(" - %s \n", match)
			// Ask the user if they want to continue or exit
		}

		// Initialize CommitSelection to track what files will be committed
		commitSelection := ship.NewCommitSelection()

		// Verify if they want to add ALL files or not
		addFilesPrompt := promptui.Select{
			Label: "Select an Option",
			Items: []string{"Select ALL files (includes flagged files)", "Select Certain Files", "Exit"},
		}

		_, filesPromptResult, err := addFilesPrompt.Run()

		if err != nil {
			fmt.Printf("Prompt Failed - %s", err)
			return
		}

		if filesPromptResult == "Exit" {
			RED_CLI_PROMPT.Println("Have a good one!")
			return
		}

		if strings.Contains(filesPromptResult, "ALL") {
			GREEN_CLI_PROMPT.Println("Adding ALL Files!")
			commitSelection.SelectAll(shipRes)
		} else if strings.Contains(filesPromptResult, "Certain") {
			// Get all files for selection
			allFiles := shipRes.GetAllFilesList()

			if len(allFiles) == 0 {
				fmt.Println("No files to select")
				return
			}

			// Build display items with indicators for dirty files
			displayItems := make([]string, len(allFiles))
			for i, file := range allFiles {
				if shipRes.IsDirtyFile(file) {
					displayItems[i] = fmt.Sprintf("%s (flagged)", file)
				} else {
					displayItems[i] = file
				}
			}

			// Multi-select loop - keep selecting until user is done
			var selectedFiles []string
			remainingFiles := allFiles
			remainingDisplay := displayItems

			for {
				if len(remainingFiles) == 0 {
					break
				}

				// Add "Done selecting" option at the top
				selectItems := append([]string{"Done selecting"}, remainingDisplay...)

				fileSelectPrompt := promptui.Select{
					Label: fmt.Sprintf("Select files to add (%d selected)", len(selectedFiles)),
					Items: selectItems,
				}

				idx, _, err := fileSelectPrompt.Run()
				if err != nil {
					fmt.Printf("Prompt Failed - %s", err)
					return
				}

				// If user selected "Done selecting"
				if idx == 0 {
					break
				}

				// Adjust index for the "Done selecting" option
				actualIdx := idx - 1
				selectedFile := remainingFiles[actualIdx]
				selectedFiles = append(selectedFiles, selectedFile)

				GREEN_CLI_PROMPT.Printf("Added: %s\n", selectedFile)

				// Remove selected file from remaining lists
				remainingFiles = append(remainingFiles[:actualIdx], remainingFiles[actualIdx+1:]...)
				remainingDisplay = append(remainingDisplay[:actualIdx], remainingDisplay[actualIdx+1:]...)
			}

			if len(selectedFiles) == 0 {
				RED_CLI_PROMPT.Println("No files selected. Exiting.")
				return
			}

			commitSelection.SelectFiles(shipRes, selectedFiles)
			GREEN_CLI_PROMPT.Printf("Selected %d file(s) for commit\n", len(selectedFiles))
		}

		// Display summary of selected files
		if len(commitSelection.SelectedFiles) > 0 {
			fmt.Println("\nFiles to be added:")
			for file, status := range commitSelection.SelectedFiles {
				if shipRes.IsDirtyFile(file) {
					RED_CLI_PROMPT.Printf("  %s", file)
					WHITE_UNDERLINE.Printf(" %s\n", status)
				} else {
					GREEN_CLI_PROMPT.Printf("  %s\n", file)
				}
			}

			if commitSelection.IncludesDirty {
				RED_CLI_PROMPT.Println("\n Warning: Selection includes flagged files!")
			}
		}

		// Phase 4: Use commitSelection.SelectedFiles for git add operations
		csAddGit, err := commitSelection.AddGitFiles()
		GREEN_CLI_PROMPT.Println("\n", csAddGit)

		// Phase 5: Ask for a git commit message. Validate the user input
		// Maybe prevent swear words and etc
		// I think we can just store it in a variable for now
		// since we are just passing it to a func
		// Idk if we need it in the commitSelection struct
		// error_log

		commitMessagePrompt := promptui.Prompt{
			Label:     "Please Enter a Commit Message",
			Default:   "Broke everything",
			AllowEdit: true,
		}

		commitResult, err := commitMessagePrompt.Run()

		if err != nil {
			return
		}

		commitSelection.AddCommitMessage(commitResult)

		// Phase 6: Push to branchname.
		_, err = ship.PushGit(shipRes.Branchname)

		if err != nil {
			RED_CLI_PROMPT.Println(err)
		}

	case "sync":
		// Check branch
		// ask if changes want to be stashed?
		// then ask which branch they want merged in
		// fetch it from origin
		// then merge it
		fmt.Println("Sync")
	case "oops":
		fmt.Println("Oopsie. Time to revert")
	case "help":
		fmt.Println("Available commands:")
		for cmd, desc := range commands {
			fmt.Printf("  %-10s %s\n", cmd, desc)
		}
	case "quote":
		fmt.Println("If a program is slow it might have a loop in it.")
	}

}
