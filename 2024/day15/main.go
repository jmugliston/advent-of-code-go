package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/jmugliston/aoc/grid"
	"github.com/jmugliston/aoc/parsing"
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

func ScaleUp(g grid.StringGrid) grid.StringGrid {
	newRows := make([][]string, 0)

	for _, row := range g {
		newRow := make([]string, len(row)*2)

		for i := 0; i < len(row); i++ {
			if row[i] == "O" {
				newRow[2*i] = "["
				newRow[(2*i)+1] = "]"
			} else if row[i] == "@" {
				newRow[2*i] = "@"
				newRow[(2*i)+1] = "."
			} else {
				newRow[2*i] = row[i]
				newRow[(2*i)+1] = row[i]
			}
		}

		newRows = append(newRows, newRow)
	}

	newGrid := grid.StringGrid(newRows)

	return newGrid
}

func parseInput(input string, double bool) (grid.StringGrid, []grid.Direction) {
	split := strings.Split(input, "\n\n")

	warehouseMap := grid.Parse(split[0])

	if double {
		warehouseMap = ScaleUp(warehouseMap)
	}

	directions := make([]grid.Direction, 0)
	for _, line := range parsing.ReadLines(split[1]) {
		chars := parsing.ReadCharacters(line)
		for _, char := range chars {
			directions = append(directions, grid.ParseDirection(char))
		}
	}

	return warehouseMap, directions
}

func robotMove(
	warehouseMap *grid.StringGrid,
	direction grid.Direction) *grid.StringGrid {

	robotPosition := warehouseMap.Find("@")
	nextPosition := robotPosition.NextPoint(direction)

	for {
		if warehouseMap.GetPoint(nextPosition) == "#" {
			break
		}
		if warehouseMap.GetPoint(nextPosition) == "." {
			warehouseMap.SetPoint(nextPosition, "@")
			warehouseMap.SetPoint(robotPosition, ".")
			break
		} else {
			nextBoxPosition := nextPosition
			for {
				nextBoxPosition = nextBoxPosition.NextPoint(direction)
				if warehouseMap.GetPoint(nextBoxPosition) == "#" {
					break
				}
				if warehouseMap.GetPoint(nextBoxPosition) == "." {
					// Move all the boxes forward
					warehouseMap.SetPoint(nextBoxPosition, "O")
					warehouseMap.SetPoint(nextPosition, "@")
					warehouseMap.SetPoint(robotPosition, ".")
					break
				}
			}

			break
		}
	}

	return warehouseMap
}

// Detect if a scaled box can move North / South
func canBoxMoveNorthSouth(warehouseMap *grid.StringGrid, direction grid.Direction, boxPositions []grid.Point, allBoxPositions [][]grid.Point) (bool, [][]grid.Point) {
	leftPosition := boxPositions[0].NextPoint(direction)
	rightPosition := boxPositions[1].NextPoint(direction)

	left := warehouseMap.GetPoint(leftPosition)
	right := warehouseMap.GetPoint(rightPosition)

	if left == "#" || right == "#" {
		return false, allBoxPositions
	}

	if left == "." && right == "." {
		return true, append(allBoxPositions, boxPositions)
	}

	allBoxPositions = append(allBoxPositions, boxPositions)

	switch {
	case left == "." && right == "[":
		return canBoxMoveNorthSouth(warehouseMap, direction, []grid.Point{rightPosition, rightPosition.NextPoint(grid.East)}, allBoxPositions)
	case left == "]" && right == ".":
		return canBoxMoveNorthSouth(warehouseMap, direction, []grid.Point{leftPosition.NextPoint(grid.West), leftPosition}, allBoxPositions)
	case left == "[" && right == "]":
		return canBoxMoveNorthSouth(warehouseMap, direction, []grid.Point{leftPosition, rightPosition}, allBoxPositions)
	case left == "]" && right == "[":
		rCanMove, rBoxes := canBoxMoveNorthSouth(warehouseMap, direction, []grid.Point{rightPosition, rightPosition.NextPoint(grid.East)}, [][]grid.Point{})
		lCanMove, lBoxes := canBoxMoveNorthSouth(warehouseMap, direction, []grid.Point{leftPosition.NextPoint(grid.West), leftPosition}, [][]grid.Point{})
		if lCanMove && rCanMove {
			return true, append(allBoxPositions, append(lBoxes, rBoxes...)...)
		}
	}

	return false, allBoxPositions
}

// This is horrible code but it works...
func robotMoveScaled(warehouseMap *grid.StringGrid, direction grid.Direction) *grid.StringGrid {
	robotPosition := warehouseMap.Find("@")
	nextPosition := robotPosition.NextPoint(direction)
	next := warehouseMap.GetPoint(nextPosition)

	if next == "#" {
		return warehouseMap
	}

	if next == "." {
		warehouseMap.SetPoint(nextPosition, "@")
		warehouseMap.SetPoint(robotPosition, ".")
		return warehouseMap
	}

	if next == "[" || next == "]" {
		if direction == grid.North || direction == grid.South {
			var canMove bool
			var positions [][]grid.Point
			var box []grid.Point

			if next == "[" {
				box = []grid.Point{nextPosition, nextPosition.NextPoint(grid.East)}
			} else if next == "]" {
				box = []grid.Point{nextPosition.NextPoint(grid.West), nextPosition}
			}

			canMove, positions = canBoxMoveNorthSouth(warehouseMap, direction, box, [][]grid.Point{})

			if canMove {
				slices.SortFunc(positions, func(a, b []grid.Point) int {
					if direction == grid.North {
						return a[0].Y - b[0].Y
					}
					return b[0].Y - a[0].Y
				})

				for _, box := range positions {
					warehouseMap.SetPoint(box[0].NextPoint(direction), "[")
					warehouseMap.SetPoint(box[1].NextPoint(direction), "]")
					warehouseMap.SetPoint(box[0], ".")
					warehouseMap.SetPoint(box[1], ".")
				}

				warehouseMap.SetPoint(nextPosition, "@")
				warehouseMap.SetPoint(robotPosition, ".")
			}
		} else {
			// Find the next point that is not a box
			count := 1
			nextPoint := nextPosition.NextPoint(direction)
			for {
				if warehouseMap.GetPoint(nextPoint) == "#" {
					break
				}
				if warehouseMap.GetPoint(nextPoint) == "." {
					// Work backwards and move the boxes along 1 space
					previousPoint := nextPoint
					for count > 0 {
						warehouseMap.SetPoint(previousPoint, warehouseMap.GetPoint(previousPoint.NextPoint(direction.Opposite())))
						previousPoint = previousPoint.NextPoint(direction.Opposite())
						count -= 1
					}
					warehouseMap.SetPoint(previousPoint, "@")
					warehouseMap.SetPoint(robotPosition, ".")
					break
				}
				count += 1
				nextPoint = nextPoint.NextPoint(direction)
			}
		}
	}

	return warehouseMap
}

func Part1(input string) int {
	warehouseMap, directions := parseInput(input, false)

	for _, direction := range directions {
		robotMove(&warehouseMap, direction)
	}

	boxLocations := warehouseMap.FindAll("O")

	result := 0
	for _, box := range boxLocations {
		result += box.X + (box.Y * 100)
	}

	return result
}

func Part2(input string) int {
	warehouseMap, directions := parseInput(input, true)

	for _, direction := range directions {
		robotMoveScaled(&warehouseMap, direction)
	}

	boxLocations := warehouseMap.FindAll("[")

	result := 0
	for _, box := range boxLocations {
		result += box.X + (box.Y * 100)
	}

	return result
}
