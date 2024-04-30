package main

import (
	"fmt"
	"math"
	"os"
	"slices"
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
	fmt.Println(part1Answer)

	part2Answer := Part2(inputString)
	fmt.Println(part2Answer)
}

type coord struct {
	x int
	y int
}

type partNumber struct {
	number string
	coords []coord
}

func mapSymbols(grid [][]string, isGear bool) []coord {

	var symbolCoords []coord

	ignoreChars := []string{".", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for y := range grid {
		for x, char := range grid[y] {
			if !slices.Contains(ignoreChars, char) {
				if !isGear || (isGear && char == "*") {
					symbolCoords = append(symbolCoords, coord{x, y})
				}
			}
		}
	}

	return symbolCoords
}

func mapParts(grid [][]string) []partNumber {
	var partNumbers []partNumber

	digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for y := range grid {

		var currentPart partNumber

		for x, char := range grid[y] {
			fmt.Println(x, char)
			if slices.Contains(digits, char) {
				currentPart.number += char
				currentPart.coords = append(currentPart.coords, coord{x, y})
			} else {
				if currentPart.number != "" {
					fmt.Println("Adding part", currentPart.number)
					partNumbers = append(partNumbers, currentPart)
					currentPart = partNumber{}
				}
			}
		}

		if currentPart.number != "" {
			partNumbers = append(partNumbers, currentPart)
		}
	}

	return partNumbers
}

func isAdjacent(a coord, b coord) bool {
	dx := math.Abs(float64(a.x) - float64(b.x))
	dy := math.Abs(float64(a.y) - float64(b.y))
	return dx <= 1 && dy <= 1
}

func getAdjacentNumbers(gear coord, numbers []partNumber) []partNumber {
	var adjacent []partNumber
	for _, number := range numbers {
		for _, coord := range number.coords {
			if isAdjacent(gear, coord) {
				adjacent = append(adjacent, number)
				break
			}
		}
	}
	return adjacent
}

func Part1(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	grid := make([][]string, len(lines))

	for i := range lines {
		grid[i] = strings.Split(lines[i], "")
	}

	symbols := mapSymbols(grid, false)
	parts := mapParts(grid)

	fmt.Println(parts)

	total := 0

	// For each gear, check for overlap with each number coords
	for _, symbol := range symbols {
		adjacentNumbers := getAdjacentNumbers(symbol, parts)
		// Add all number that have an adjacent part
		if len(adjacentNumbers) > 0 {
			for _, number := range adjacentNumbers {
				n, _ := strconv.Atoi(number.number)
				total += n
			}
		}
	}

	return total
}

func Part2(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	grid := make([][]string, len(lines))

	for i := range lines {
		grid[i] = strings.Split(lines[i], "")
	}

	gears := mapSymbols(grid, true)
	parts := mapParts(grid)

	total := 0

	// For each gear, check for overlap with each number coords
	for _, gear := range gears {
		adjacentNumbers := getAdjacentNumbers(gear, parts)

		// We're only interested in gears that have exactly 2 adjacent numbers
		if len(adjacentNumbers) == 2 {
			a, _ := strconv.Atoi(adjacentNumbers[0].number)
			b, _ := strconv.Atoi(adjacentNumbers[1].number)
			total += a * b
		}
	}

	return total
}
