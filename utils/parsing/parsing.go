package parsing

import (
	"regexp"
	"strconv"
	"strings"
)

func ReadNumbers(input string) []int {
	var numbers []int

	numbersString := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(input), -1)

	for _, char := range numbersString {
		number, _ := strconv.Atoi(string(char))
		numbers = append(numbers, number)
	}

	return numbers
}

func ReadLines(input string) []string {
	return strings.Split(strings.TrimSpace(input), "\n")
}

func ReadLinesOfNumbers(input string) [][]int {
	var lines [][]int

	for _, line := range ReadLines(input) {
		lines = append(lines, ReadNumbers(line))
	}

	return lines
}
