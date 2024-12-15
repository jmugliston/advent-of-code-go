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

	nextPosition := robotPosition.GetNextPointInDirection(direction)
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
				nextBoxPosition = nextBoxPosition.GetNextPointInDirection(direction)
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

func canBoxMove(warehouseMap *grid.StringGrid, direction grid.Direction, boxPositions []grid.Point, allboxPositions [][]grid.Point) (bool, [][]grid.Point) {
	a := boxPositions[0].GetNextPointInDirection(direction)
	b := boxPositions[1].GetNextPointInDirection(direction)

	aValue := warehouseMap.GetPoint(a)
	bValue := warehouseMap.GetPoint(b)

	if aValue == "#" || bValue == "#" {
		return false, allboxPositions
	}

	if aValue == "." && bValue == "." {
		return true, append(allboxPositions, []grid.Point{boxPositions[0], boxPositions[1]})
	}

	allboxPositions = append(allboxPositions, []grid.Point{boxPositions[0], boxPositions[1]})

	if aValue == "[" && bValue == "." {
		return canBoxMove(warehouseMap, direction, []grid.Point{a, a.GetNextPointInDirection(grid.West)}, allboxPositions)
	}

	if aValue == "." && bValue == "[" {
		return canBoxMove(warehouseMap, direction, []grid.Point{b, b.GetNextPointInDirection(grid.East)}, allboxPositions)
	}

	if aValue == "." && bValue == "]" {
		return canBoxMove(warehouseMap, direction, []grid.Point{b.GetNextPointInDirection(grid.West), b}, allboxPositions)
	}

	if aValue == "]" && bValue == "." {
		return canBoxMove(warehouseMap, direction, []grid.Point{a.GetNextPointInDirection(grid.West), a}, allboxPositions)
	}

	if aValue == "[" && bValue == "]" {
		return canBoxMove(warehouseMap, direction, []grid.Point{a, b}, allboxPositions)
	}

	if aValue == "]" && bValue == "[" {
		rCanMove, rBoxes := canBoxMove(warehouseMap, direction, []grid.Point{b, b.GetNextPointInDirection(grid.East)}, [][]grid.Point{})
		lCanMove, lBoxes := canBoxMove(warehouseMap, direction, []grid.Point{a.GetNextPointInDirection(grid.West), a}, [][]grid.Point{})
		if lCanMove && rCanMove {
			return true, append(allboxPositions, append(lBoxes, rBoxes...)...)
		} else {
			return false, allboxPositions
		}
	}

	return false, allboxPositions
}

func sortYAscending(a, b []grid.Point) int {
	if a[0].Y < b[0].Y {
		return -1
	} else if a[0].Y > b[0].Y {
		return 1
	}
	return 0
}

func sortYDescending(a, b []grid.Point) int {
	if a[0].Y < b[0].Y {
		return 1
	} else if a[0].Y > b[0].Y {
		return -1
	}
	return 0
}

func robotMoveScaled(
	warehouseMap *grid.StringGrid,
	direction grid.Direction) *grid.StringGrid {

	robotPosition := warehouseMap.Find("@")

	nextPosition := robotPosition.GetNextPointInDirection(direction)

	if warehouseMap.GetPoint(nextPosition) == "#" {
		// Do nothing
		return warehouseMap
	}

	if warehouseMap.GetPoint(nextPosition) == "." {
		// Just move into that space
		warehouseMap.SetPoint(nextPosition, "@")
		warehouseMap.SetPoint(robotPosition, ".")
		return warehouseMap
	}

	if warehouseMap.GetPoint(nextPosition) == "[" || warehouseMap.GetPoint(nextPosition) == "]" {
		if direction == grid.North || direction == grid.South {
			if warehouseMap.GetPoint(nextPosition) == "#" {
				// Hit a wall - do nothing
			}

			if warehouseMap.GetPoint(nextPosition) == "[" {
				boxPositions := []grid.Point{nextPosition, nextPosition.GetNextPointInDirection(grid.East)}
				canMove, positions := canBoxMove(warehouseMap, direction, boxPositions, [][]grid.Point{})
				if canMove {
					if direction == grid.North {
						slices.SortFunc(positions, sortYAscending)
					} else {
						slices.SortFunc(positions, sortYDescending)
					}
					for _, box := range positions {
						warehouseMap.SetPoint(box[0].GetNextPointInDirection(direction), "[")
						warehouseMap.SetPoint(box[1].GetNextPointInDirection(direction), "]")
						warehouseMap.SetPoint(box[0], ".")
						warehouseMap.SetPoint(box[1], ".")
					}
					warehouseMap.SetPoint(robotPosition.GetNextPointInDirection(direction), "@")
					warehouseMap.SetPoint(robotPosition.GetNextPointInDirection(direction).GetNextPointInDirection(grid.East), ".")
					warehouseMap.SetPoint(robotPosition, ".")
				}
			}

			if warehouseMap.GetPoint(nextPosition) == "]" {
				boxPositions := []grid.Point{nextPosition.GetNextPointInDirection(grid.West), nextPosition}
				canMove, positions := canBoxMove(warehouseMap, direction, boxPositions, [][]grid.Point{})
				if canMove {
					if direction == grid.North {
						slices.SortFunc(positions, sortYAscending)
					} else {
						slices.SortFunc(positions, sortYDescending)
					}

					for _, box := range positions {
						warehouseMap.SetPoint(box[0].GetNextPointInDirection(direction), "[")
						warehouseMap.SetPoint(box[1].GetNextPointInDirection(direction), "]")
						warehouseMap.SetPoint(box[0], ".")
						warehouseMap.SetPoint(box[1], ".")
					}
					warehouseMap.SetPoint(robotPosition.GetNextPointInDirection(direction), "@")
					warehouseMap.SetPoint(robotPosition.GetNextPointInDirection(direction).GetNextPointInDirection(grid.West), ".")
					warehouseMap.SetPoint(robotPosition, ".")
				}
			}
		} else {
			// Find the next point that is not a box
			count := 1
			nextPoint := nextPosition.GetNextPointInDirection(direction)
			for {
				if warehouseMap.GetPoint(nextPoint) == "#" {
					// Hit a wall - do nothing
					break
				}
				if warehouseMap.GetPoint(nextPoint) == "." {
					// Work backwards and move the boxes along 1 space
					previousPoint := nextPoint
					for count > 0 {
						warehouseMap.SetPoint(previousPoint, warehouseMap.GetPoint(previousPoint.GetNextPointInDirection(grid.GetOppositeDirection(direction))))
						previousPoint = previousPoint.GetNextPointInDirection(grid.GetOppositeDirection(direction))
						count -= 1
					}
					warehouseMap.SetPoint(previousPoint, "@")
					warehouseMap.SetPoint(robotPosition, ".")
					break
				}
				count += 1
				nextPoint = nextPoint.GetNextPointInDirection(direction)
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
