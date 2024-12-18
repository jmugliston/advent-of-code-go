package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	lines := strings.Split(input, "\n")

	total := 0

	for _, line := range lines {
		total += sumFirstAndLastDigits(line)
	}

	return total
}

func Part2(input string) int {
	lines := strings.Split(input, "\n")

	numberMap := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
		"zero":  "0",
	}

	total := 0

	for _, line := range lines {

		newLine := strings.Clone(line)

		// Remap the line e.g. two1nine becomes two2two1nine9nine
		// This helps solve for overlapping values
		for key, val := range numberMap {
			newLine = strings.ReplaceAll(newLine, key, key+val+key)
		}

		total += sumFirstAndLastDigits(newLine)
	}

	return total

}

func sumFirstAndLastDigits(line string) int {

	var firstNumber string
	var lastNumber string

	for _, char := range line {
		if char >= '0' && char <= '9' {
			if firstNumber == "" {
				firstNumber = string(char)
			}
			lastNumber = string(char)
		}
	}

	num, _ := strconv.Atoi(firstNumber + lastNumber)

	return num
}
