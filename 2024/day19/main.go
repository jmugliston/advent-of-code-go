package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jmugliston/aoc/parsing"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")
var exampleFlag = flag.Bool("example", false, "Use the example instead of the puzzle input")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)

	inputFile := "input.txt"
	if *exampleFlag {
		inputFile = "example.txt"
	}
	path := filepath.Join(filepath.Dir(filename), "input", inputFile)

	input, err := os.ReadFile(path)

	if err != nil {
		panic("Could not find the input file")
	}

	if *partFlag == "1" {
		fmt.Println(Part1(string(input)))
	} else {
		fmt.Println(Part2(string(input)))
	}
}

func parseInput(input string) ([]string, []string) {
	split := strings.Split(input, "\n\n")

	patterns := strings.Split(split[0], ", ")
	designs := parsing.ReadLines(split[1])

	return patterns, designs
}

// Return a count of the number of ways the design can be matched
func matchCount(design string, patterns []string, cache *map[string]int) int {
	if len(design) == 0 {
		return 1
	}

	if _, ok := (*cache)[design]; ok {
		return (*cache)[design]
	}

	count := 0
	for _, pattern := range patterns {
		if strings.HasPrefix(design, pattern) {
			count += matchCount(design[len(pattern):], patterns, cache)
		}
	}

	(*cache)[design] = count

	return count
}

func Part1(input string) int {
	patterns, designs := parseInput(input)

	cache := make(map[string]int)

	count := 0
	for _, design := range designs {
		if matchCount(design, patterns, &cache) > 0 {
			count++
		}
	}

	return count
}

func Part2(input string) int {
	patterns, designs := parseInput(input)

	cache := make(map[string]int)

	count := 0
	for _, design := range designs {
		count += matchCount(design, patterns, &cache)
	}

	return count
}
