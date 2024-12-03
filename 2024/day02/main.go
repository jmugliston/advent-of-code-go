package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/atheius/aoc/parsing"
	"github.com/atheius/aoc/utils"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "input", "input.txt")

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

func isDiffSafe(a int, b int, isIncreasing bool) bool {
	abs := utils.Abs(b - a)
	if abs < 1 || abs > 3 {
		return false
	}
	if (isIncreasing && b < a) || (!isIncreasing && b > a) {
		return false
	}
	return true
}

func isReportSafe(line []int) bool {
	isIncreasing := line[1] > line[0]
	for i := 1; i < len(line); i++ {
		if !isDiffSafe(line[i-1], line[i], isIncreasing) {
			return false
		}
	}
	return true
}

func Part1(input string) int {
	lines := parsing.ReadLinesOfNumbers(input)

	numSafeReports := 0

	for _, line := range lines {
		isSafe := isReportSafe(line)
		if isSafe {
			numSafeReports++
		}
	}

	return numSafeReports
}

func Part2(input string) int {
	lines := parsing.ReadLinesOfNumbers(input)

	numSafeReports := 0

	for _, line := range lines {
		isSafe := isReportSafe(line)

		if !isSafe {
			// Brute force - Try removing 1 element at a time
			for i := 0; i < len(line); i++ {
				next := utils.RemoveIndex(line, i)
				isSafe = isReportSafe(next)
				if isSafe {
					break
				}
			}
		}

		if isSafe {
			numSafeReports++
		}
	}

	return numSafeReports
}
