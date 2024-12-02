package parsing

import (
	"regexp"
	"strconv"
	"strings"
)

func ReadDigits(input string) []int {
	var digits []int

	numbersString := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(input), -1)

	for _, char := range numbersString {
		digit, _ := strconv.Atoi(string(char))
		digits = append(digits, digit)
	}

	return digits
}

func ReadLines(input string) []string {
	return strings.Split(strings.TrimSpace(input), "\n")
}

func ReadLinesOfDigits(input string) [][]int {
	var lines [][]int

	for _, line := range ReadLines(input) {
		lines = append(lines, ReadDigits(line))
	}

	return lines
}
