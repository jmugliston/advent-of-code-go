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
