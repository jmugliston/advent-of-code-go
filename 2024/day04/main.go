package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jmugliston/aoc/grid"
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
	wordsearch := grid.Parse(input)

	total := 0
	for _, point := range wordsearch.FindAll("X") {
		for _, direction := range grid.Directions {
			nextNPoints := grid.NextPoints(point.AddDirection(direction), 3)
			wordPoints := append([]grid.Point{point}, nextNPoints...)
			word := strings.Join(wordsearch.GetPoints(wordPoints)[:], "")
			if word == "XMAS" {
				total += 1
			}
		}
	}

	return total
}

func Part2(input string) int {
	wordsearch := grid.Parse(input)

	total := 0
	for _, point := range wordsearch.FindAll("A") {
		points := point.Neighbours()

		nw := wordsearch.GetPoint(points.NorthWest)
		se := wordsearch.GetPoint(points.SouthEast)
		ne := wordsearch.GetPoint(points.NorthEast)
		sw := wordsearch.GetPoint(points.SouthWest)

		// NorthWest -> SouthEast diagonal
		if nw == "M" && se == "S" || se == "M" && nw == "S" {
			// NorthEast -> SouthWest diagonal
			if ne == "M" && sw == "S" || sw == "M" && ne == "S" {
				total += 1
			}
		}
	}

	return total
}
