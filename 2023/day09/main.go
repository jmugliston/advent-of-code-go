package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/atheius/aoc/parsing"
	"github.com/atheius/aoc/utils"
)

func main() {
	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the input file!")
	}

	inputString := string(input)

	part1Answer := Part1(inputString)
	fmt.Println(part1Answer)

	part2Answer := Part2(inputString)
	fmt.Println(part2Answer)
}

func allZeros(nums []int) bool {
	allZeros := true
	for _, num := range nums {
		if num != 0 {
			allZeros = false
			break
		}
	}
	return allZeros
}

func getNextLines(line []int) [][]int {
	nextLines := [][]int{line}

	for {
		var nextLine []int

		lastLine := nextLines[len(nextLines)-1]
		for i := 0; i < len(lastLine)-1; i++ {
			nextLine = append(nextLine, lastLine[i+1]-lastLine[i])
		}

		if allZeros(nextLine) {
			break
		}

		nextLines = append(nextLines, nextLine)

	}

	return nextLines
}

func getNextNumber(lines [][]int, part2 bool) int {
	slices.Reverse(lines)

	nextNumber := 0
	for _, line := range lines {
		if part2 {
			nextNumber = line[0] - nextNumber
		} else {
			nextNumber = nextNumber + line[len(line)-1]
		}
	}

	return nextNumber
}

func Part1(input string) int {

	lines := strings.Split(strings.TrimSpace(input), "\n")

	var numbers []int
	for _, line := range lines {
		nums := parsing.ReadDigits(line)
		nextLines := getNextLines(nums)
		nextNumber := getNextNumber(nextLines, false)
		numbers = append(numbers, nextNumber)
	}

	return utils.Sum(numbers)
}

func Part2(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var numbers []int
	for _, line := range lines {
		nums := parsing.ReadDigits(line)
		nextLines := getNextLines(nums)
		nextNumber := getNextNumber(nextLines, true)
		numbers = append(numbers, nextNumber)
	}

	return utils.Sum(numbers)
}
