package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/atheius/aoc/parsing"
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

func getRangeMaps(lines []string) [][][]int {
	rangeMaps := [][][]int{}
	for _, rangeMap := range lines {
		rangeLines := strings.Split(rangeMap, "\n")[1:]
		var ranges [][]int
		for _, nextRange := range rangeLines {
			ranges = append(ranges, parsing.ReadDigits(nextRange))
		}
		rangeMaps = append(rangeMaps, ranges)
	}
	return rangeMaps
}

func seedToLocation(seed int, rangeMaps [][][]int) int {
	currentValue := seed
	for _, nextRangeMap := range rangeMaps {
		for _, nextRange := range nextRangeMap {
			if currentValue >= nextRange[1] && currentValue < nextRange[1]+nextRange[2] {
				currentValue = nextRange[0] + (currentValue - nextRange[1])
				break
			}
		}
	}
	return currentValue
}

func locationToSeed(location int, rangeMaps [][][]int) int {
	currentValue := location
	// Need to go backwards through the ranges
	for i := len(rangeMaps) - 1; i >= 0; i-- {
		nextRangeMap := rangeMaps[i]
		for _, nextRange := range nextRangeMap {
			if currentValue >= nextRange[0] && currentValue < nextRange[0]+nextRange[2] {
				currentValue = nextRange[1] + (currentValue - nextRange[0])
				break
			}
		}
	}
	return currentValue
}

func Part1(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n\n")

	seedLine := strings.Split(lines[0], ": ")
	seeds := parsing.ReadDigits(seedLine[1])

	rangeMaps := getRangeMaps(lines[1:])

	seedLocations := []int{}
	for _, seed := range seeds {
		seedLocations = append(seedLocations, seedToLocation(seed, rangeMaps))
	}

	return slices.Min(seedLocations)
}

func Part2(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n\n")

	seedLine := strings.Split(lines[0], ": ")
	seedRangeList := parsing.ReadDigits(seedLine[1])

	seedRanges := [][]int{}

	for i := 0; i <= len(seedRangeList); i = i + 2 {
		if (i + 1) < len(seedRangeList) {
			seedRanges = append(seedRanges, []int{seedRangeList[i], seedRangeList[i+1]})
		}
	}

	rangeMaps := getRangeMaps(lines[1:])

	// Brute force from location 1 to find the the first valid seed

	location := 1
	for {

		seed := locationToSeed(location, rangeMaps)

		for _, seedRange := range seedRanges {
			if seed >= seedRange[0] && seed < seedRange[0]+seedRange[1] {
				return location
			}
		}

		location += 1

	}

}
