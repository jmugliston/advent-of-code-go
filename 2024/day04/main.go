package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func checkForWord(wordsearch *grid.StringGrid, word string, startPoint grid.Point, direction grid.Direction) bool {
	nextNPoints := grid.GetNextNPointsInDirection(grid.PointWithDirection{
		X:         startPoint.X,
		Y:         startPoint.Y,
		Direction: direction,
	}, len(word)-1)

	points := append([]grid.Point{startPoint}, nextNPoints...)

	for idx, point := range points {
		if !wordsearch.IsPointInGrid(point) || string(word[idx]) != (*wordsearch)[point.Y][point.X] {
			return false
		}
	}

	return true
}

func Part1(input string) int {
	wordsearch := grid.Parse(input)
	searchWord := "XMAS"

	total := 0
	for y, row := range wordsearch {
		for x, cell := range row {
			if cell == "X" {
				for _, direction := range grid.Directions {
					// Check for the word XMAS in all directions from this "X"
					if checkForWord(&wordsearch, searchWord, grid.Point{X: x, Y: y}, direction) {
						total += 1
					}
				}
			}
		}
	}

	return total
}

func Part2(input string) int {
	wordsearch := grid.Parse(input)
	searchWord := "MAS"

	total := 0
	for y, row := range wordsearch {
		for x, cell := range row {
			if cell == "A" {
				points := grid.SurroundingPoints(grid.Point{X: x, Y: y})
				// Check both diagonals for the word "MAS" from this "A"
				diagonalA := checkForWord(&wordsearch, searchWord, points.NorthWest, grid.SouthEast) || checkForWord(&wordsearch, searchWord, points.SouthEast, grid.NorthWest)
				diagonalB := checkForWord(&wordsearch, searchWord, points.NorthEast, grid.SouthWest) || checkForWord(&wordsearch, searchWord, points.SouthWest, grid.NorthEast)
				if diagonalA && diagonalB {
					total += 1
				}
			}
		}
	}

	return total
}
