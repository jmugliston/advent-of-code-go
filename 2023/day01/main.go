package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the input file!")
	}

	inputString := string(input)

	part1Answer := Part1(inputString)
	fmt.Printf("Part 1: %v\n", part1Answer)

	part2Answer := Part2(inputString)
	fmt.Printf("Part 2: %v\n", part2Answer)
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
