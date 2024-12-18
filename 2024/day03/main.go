package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

func Part1(input string) int {
	corruptedMemory := strings.TrimSpace(input)

	pattern := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := pattern.FindAllStringSubmatch(corruptedMemory, -1)

	total := 0

	for _, match := range matches {
		a, _ := strconv.Atoi(match[1])
		b, _ := strconv.Atoi(match[2])
		total += a * b
	}

	return total
}

func Part2(input string) int {
	corruptedMemory := strings.TrimSpace(input)

	pattern := regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	matches := pattern.FindAllStringSubmatch(corruptedMemory, -1)

	total := 0
	multEnabled := true

	for _, match := range matches {
		if match[0] == "do()" {
			multEnabled = true
		} else if match[0] == "don't()" {
			multEnabled = false
		} else if multEnabled {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			total += a * b
		}
	}

	return total
}
