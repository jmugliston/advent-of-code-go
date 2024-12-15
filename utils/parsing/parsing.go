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

func ReadCharacters(input string) []string {
	var characters []string

	for _, char := range input {
		characters = append(characters, string(char))
	}

	return characters
}

func ReadDigits(input string) []int {
	var digits []int

	for _, char := range input {
		digit, _ := strconv.Atoi(string(char))
		digits = append(digits, digit)
	}

	return digits
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
