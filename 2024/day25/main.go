package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
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
	}
}

func parseInput(input string) ([][]int, [][]int) {
	grids := strings.Split(strings.TrimSpace(input), "\n\n")

	locks := [][]int{}
	keys := [][]int{}

	for _, grid := range grids {
		lines := strings.Split(grid, "\n")
		isLock := lines[0] == "#####"

		if !isLock {
			slices.Reverse(lines)
		}

		counts := make([]int, len(lines[0]))
		for _, line := range lines[1:] {
			for idx, char := range line {
				if char == '#' {
					counts[idx]++
				}
			}
		}

		if isLock {
			locks = append(locks, counts)
		} else {
			keys = append(keys, counts)
		}
	}

	return locks, keys
}

func Part1(input string) int {
	locks, keys := parseInput(input)

	valid := 0
	for _, lock := range locks {
		for _, key := range keys {
			isValid := true
			for i := 0; i < 5; i++ {
				if lock[i]+key[i] >= 6 {
					isValid = false
					break
				}
			}
			if isValid {
				valid++
			}
		}
	}

	return valid
}
