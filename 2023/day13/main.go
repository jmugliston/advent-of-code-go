package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atheius/aoc/grid"
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

func getGrids(input string) []grid.StringGrid {
	split := strings.Split(input, "\n\n")

	var grids []grid.StringGrid

	for _, line := range split {
		grids = append(grids, grid.Parse(line))
	}

	return grids
}

func getReflectionLineAlt(grid [][]string, isSmudged bool) int {

	for y := 0; y < len(grid); y++ {

		isReflectionLine := false
		smudgeCount := 0

		for offset := 0; offset < len(grid); offset++ {

			if y-offset < 0 || y+offset+1 >= len(grid) {
				// Out of bounds
				break
			}

			for x := 0; x < len(grid[y]); x++ {
				if grid[y-offset][x] != grid[y+offset+1][x] {
					smudgeCount += 1
				}
			}

			if (!isSmudged && smudgeCount > 0) || (isSmudged && smudgeCount > 1) {
				break
			} else {
				isReflectionLine = true
			}
		}

		if isReflectionLine {
			if (!isSmudged && smudgeCount == 0) || (isSmudged && smudgeCount == 1) {
				return y
			}
		}
	}

	return -1
}

func Part1(input string) int {

	grids := getGrids(input)

	result := 0
	for _, nextGrid := range grids {

		horizontalReflectionLine := getReflectionLineAlt(nextGrid, false)

		if horizontalReflectionLine != -1 {
			result += 100 * (horizontalReflectionLine + 1)
		} else {
			transposedGrid := nextGrid.Transpose()
			verticalReflectionLine := getReflectionLineAlt(transposedGrid, false)
			result += verticalReflectionLine + 1
		}
	}

	return result
}

func Part2(input string) int {
	grids := getGrids(input)

	result := 0
	for _, nextGrid := range grids {

		horizontalReflectionLine := getReflectionLineAlt(nextGrid, true)

		if horizontalReflectionLine != -1 {
			result += 100 * (horizontalReflectionLine + 1)
		} else {
			transposedGrid := nextGrid.Transpose()
			verticalReflectionLine := getReflectionLineAlt(transposedGrid, true)
			result += verticalReflectionLine + 1
		}
	}

	return result
}
