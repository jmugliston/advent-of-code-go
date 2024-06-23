package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atheius/aoc/grid"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "input", "input.txt")

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
