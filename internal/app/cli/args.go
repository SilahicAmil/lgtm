package cli

import "errors"

func Parse(args []string) (string, *Options, error) {
	if len(args) == 0 {
		return "", nil, errors.New("command required")
	}

	cmd := args[0]

	flagArgs := args[1:]
	opts, err := ParseFlags(flagArgs)

	if err != nil {
		return "", nil, err
	}

	return cmd, opts, nil
}
