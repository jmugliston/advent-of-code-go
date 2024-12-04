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

func isWord(wordsearch grid.StringGrid, word string, startPoint grid.Point, direction grid.Direction) bool {
	next := grid.Point{X: startPoint.X, Y: startPoint.Y}

	for _, char := range word {
		if !wordsearch.IsPointInGrid(next) {
			return false
		}
		if wordsearch[next.Y][next.X] != string(char) {
			return false
		}
		next = grid.GetNextPointInDirection(grid.PointWithDirection{
			X:         next.X,
			Y:         next.Y,
			Direction: direction,
		})
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
				for _, direction := range []grid.Direction{
					grid.North,
					grid.NorthEast,
					grid.East,
					grid.SouthEast,
					grid.South,
					grid.SouthWest,
					grid.West,
					grid.NorthWest,
				} {
					// Check for the word XMAS in all directions from this "X"
					if isWord(wordsearch, searchWord, grid.Point{X: x, Y: y}, direction) {
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
				// Check both diagonals for the word "MAS"
				diagonalA := isWord(wordsearch, searchWord, points.NorthWest, grid.SouthEast) || isWord(wordsearch, searchWord, points.SouthEast, grid.NorthWest)
				diagonalB := isWord(wordsearch, searchWord, points.NorthEast, grid.SouthWest) || isWord(wordsearch, searchWord, points.SouthWest, grid.NorthEast)
				if diagonalA && diagonalB {
					total += 1
				}
			}
		}
	}

	return total
}
