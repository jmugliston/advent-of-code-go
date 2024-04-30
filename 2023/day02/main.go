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
	lines := strings.Split(strings.TrimSpace(input), "\n")

	sum := 0

	const redMax = 12
	const greenMax = 13
	const blueMax = 14

	for idx, line := range lines {

		possible := true

		game := strings.Split(line, ": ")
		turns := strings.Split(game[1], "; ")

		for _, turn := range turns {

			totals := make(map[string]int)
			colours := strings.Split(turn, ", ")

			for _, colour := range colours {
				pair := strings.Split(colour, " ")
				num, _ := strconv.Atoi(pair[0])
				colour := pair[1]
				totals[colour] += num

			}

			if totals["red"] > redMax || totals["green"] > greenMax || totals["blue"] > blueMax {
				possible = false
				break
			}
		}

		if possible {
			sum += idx + 1
		}

	}

	return sum
}

func Part2(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	power := 0

	for _, line := range lines {

		game := strings.Split(line, ": ")
		turns := strings.Split(game[1], "; ")

		max := make(map[string]int)

		for _, turn := range turns {
			colours := strings.Split(turn, ", ")
			for _, colour := range colours {
				pair := strings.Split(colour, " ")
				num, _ := strconv.Atoi(pair[0])
				colour := pair[1]
				if num > max[colour] {
					max[colour] = num
				}
			}
		}

		power += max["red"] * max["green"] * max["blue"]

	}

	return power
}
