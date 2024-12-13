package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/parsing"
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

func trimLeadingDots(row string) string {
	for i, char := range row {
		if char == '#' || char == '?' {
			return row[i:]
		}
	}
	return ""
}

func trimTrailingDots(row string) string {
	for i := len(row) - 1; i >= 0; i-- {
		if row[i] == '#' || row[i] == '?' {
			return row[:i+1]
		}
	}
	return ""
}

func getNextNumBrokenSprings(row string) int {
	springCount := 0
	for _, char := range row {
		if char == '#' {
			springCount += 1
		} else if char == '?' {
			springCount = 0
			break
		} else {
			break
		}
	}
	return springCount
}

var cache map[string]int = make(map[string]int, 1000000)

func countPossibleArrangements(row string, pattern []int) int {
	row = trimLeadingDots(row)
	row = trimTrailingDots(row)

	if row == "" {
		if len(pattern) > 0 {
			// Not a valid arrangement
			return 0
		}
		// Found a valid arrangement
		return 1
	}

	if len(pattern) == 0 {
		if strings.Contains(row, "#") {
			// Not a valid arrangement
			return 0
		}
		// Found a valid arrangement
		return 1
	}

	key := row + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(pattern)), ","), "[]")

	_, ok := cache[key]

	if ok {
		return cache[key]
	}

	arrangements := 0

	nextBroken := getNextNumBrokenSprings(row)

	if nextBroken > 0 {
		if nextBroken == pattern[0] {
			next := row[nextBroken:]
			arrangements += countPossibleArrangements(next, pattern[1:])
		}
	} else if strings.Contains(row, "?") {
		arrangements += countPossibleArrangements(strings.Replace(row, "?", ".", 1), pattern)
		arrangements += countPossibleArrangements(strings.Replace(row, "?", "#", 1), pattern)
	}

	cache[key] = arrangements

	return arrangements
}

type Record struct {
	row     string
	pattern []int
}

func parseRecords(lines []string) []Record {
	var data []Record
	for _, line := range lines {
		split := strings.Split(line, " ")
		row := split[0]
		patternString := strings.Split(split[1], ",")
		var pattern []int
		for _, char := range patternString {
			num, _ := strconv.Atoi(char)
			pattern = append(pattern, num)
		}
		data = append(data, Record{row, pattern})
	}
	return data
}

func Part1(input string) int {
	lines := parsing.ReadLines(input)

	records := parseRecords(lines)

	sum := 0
	for _, record := range records {
		arrangements := countPossibleArrangements(record.row, record.pattern)
		sum += arrangements
	}

	return sum
}

func Part2(input string) int {
	lines := parsing.ReadLines(input)

	records := parseRecords(lines)

	sum := 0
	for _, record := range records {

		var foldedRow string
		var foldedPattern []int

		for i := 0; i < 5; i++ {
			foldedRow = foldedRow + record.row
			foldedPattern = append(foldedPattern, record.pattern...)
			if i != 4 {
				foldedRow += "?"
			}
		}

		arrangements := countPossibleArrangements(foldedRow, foldedPattern)

		sum += arrangements
	}

	return sum
}
