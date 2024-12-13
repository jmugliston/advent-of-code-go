# 🎄 Advent of Code (Go) 🎄

<div align="left">
    <img
      src="https://github.com/jmugliston/advent-of-code-go/raw/HEAD/aoc-go.jpeg"
      width="150"
      height="auto"
    />
</div>

Project for AOC challenges and solutions in Go.

## Setup

Create a .env file with the following:

```
SESSION_TOKEN=<aoc-session-token>
```

Install the package locally:

```sh
make install
```

## CLI

To run the interactive cli:

```sh
aoc
```

### Usage

```
# aoc help

🎄🎄🎄 Advent of Code 🎄🎄🎄

AoC command-line tool

Usage:
  aoc [flags]
  aoc [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    Download puzzle inputs for specific year/day
  help        Help about any command
  init        Create a template folder for a specific day
  solve       Run the solution for a specific day
  submit      Submit an answer for a specific day

Flags:
  -h, --help      help for aoc
  -q, --quiet     quiet mode
  -v, --version   show version

Use "aoc [command] --help" for more information about a command
```

```
# aoc init --help

Initialise a new day

Usage:
  aoc init [flags]

Examples:
aoc init --day 1

Flags:
  -d, --day int    puzzle day (default current day of AoC event)
  -h, --help       help for init
  -y, --year int   puzzle year (default year of current or last AoC event)

Global Flags:
  -q, --quiet   quiet mode
```

```
# aoc solve --help

Run the solution for a specific day

Usage:
  aoc solve [flags]

Examples:
aoc solve --day 1 --part 1

Flags:
  -d, --day int    puzzle day (default current day of AoC event)
  -h, --help       help for solve
  -p, --part int   puzzle part (default 1)
  -y, --year int   puzzle year (default year of current or last AoC event)

Global Flags:
  -q, --quiet   quiet mode
```

```
# aoc submit --help

Submit an answer for a specific day

Usage:
  aoc submit [flags]

Examples:
aoc submit --day 1 --part 1

Flags:
  -d, --day int    puzzle day
  -h, --help       help for submit
  -p, --part int   puzzle part (default 1)
  -y, --year int   puzzle year (default year of current or last AoC event)

Global Flags:
  -q, --quiet   quiet mode
```

## Test

Run tests with:

```sh
go test ./...
```
