package cli

import "flag"

type Options struct {
	DryRun bool
	Force  bool
}

func ParseFlags(args []string) (*Options, error) {
	fs := flag.NewFlagSet("lgtm", flag.ContinueOnError)

	opts := &Options{}

	// Define flags
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Do not make changes")
	fs.BoolVar(&opts.Force, "force", false, "Force the operation")

	// Parse the args for this command
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return opts, nil
}
