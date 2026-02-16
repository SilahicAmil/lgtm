# LGTM (Looks Good To Me)

**LGTM** is a simple, fun, and productive Git CLI helper that makes staging, committing, and pushing easier - with a few fun surprises along the way.

---

## Features

- Stage and commit files quickly with the `ship` command
- Sync your branch with another branch using `sync`
- Undo last commit with `oops` if you made a mistake
- Get random quotes or motivational messages with `quote`
- Help system to show available commands

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/SilahicAmil/lgtm.git
cd lgtm
```

2. Build the binary:

```bash
make build
```

3. Run the CLI

```bash
./lgtm <command>
```

## Commands

| Command | Description                                  |
| ------- | -------------------------------------------- |
| `help`  | Show this help message                       |
| `ship`  | Stage, commit, and push files                |
| `sync`  | Sync your branch with another one            |
| `oops`  | Undo last commit and reset changes           |
| `quote` | Print a random motivational or funny message |

## Development

- Go version >= 1.21
- Uses Standard Go modules

### To run without building:

```bash
go run cmd/main.go <command>
```

## Contributing

Feel free to open issues, submit PRs, or suggest CLI commands.
