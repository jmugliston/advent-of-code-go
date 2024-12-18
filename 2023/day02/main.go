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
