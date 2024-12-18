package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"

	"github.com/jmugliston/aoc/grid"
	"github.com/jmugliston/aoc/utils"
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

func findStartPosition(maze [][]string) grid.Point {
	for y, row := range maze {
		for x, cell := range row {
			if cell == "S" {
				return grid.Point{X: x, Y: y}
			}
		}
	}
	panic("Couldn't find the start position")
}

func isValidMove(maze [][]string, previous grid.Point, current grid.Point, next grid.Point, direction grid.Direction) bool {

	if next == previous {
		// Can't move backwards
		return false
	}

	if next.Y < 0 || next.Y >= len(maze) {
		return false
	}

	if next.X < 0 || next.X >= len(maze[0]) {
		return false
	}

	currentMazeCell := maze[current.Y][current.X]
	nextMazeCell := maze[next.Y][next.X]

	if direction.String() == "North" {
		return slices.Contains([]string{"S", "|", "L", "J"}, currentMazeCell) && slices.Contains([]string{"|", "7", "F"}, nextMazeCell)
	}

	if direction.String() == "East" {
		return slices.Contains([]string{"S", "-", "F", "L"}, currentMazeCell) && slices.Contains([]string{"-", "7", "J"}, nextMazeCell)

	}

	if direction.String() == "South" {
		return slices.Contains([]string{"S", "|", "7", "F"}, currentMazeCell) && slices.Contains([]string{"|", "L", "J"}, nextMazeCell)
	}

	if direction.String() == "West" {
		return slices.Contains([]string{"S", "-", "7", "J"}, currentMazeCell) && slices.Contains([]string{"-", "L", "F"}, nextMazeCell)
	}

	panic("Invalid direction")

}

func getSteps(maze [][]string, start grid.Point) []grid.Point {

	previous := grid.Point{X: start.X, Y: start.Y}
	current := grid.Point{X: start.X, Y: start.Y}

	steps := []grid.Point{start}

	for {
		nextPoints := grid.SurroundingPoints(current)

		if isValidMove(maze, previous, current, nextPoints.North, grid.North) {
			previous = current
			current = nextPoints.North
		} else if isValidMove(maze, previous, current, nextPoints.East, grid.East) {
			previous = current
			current = nextPoints.East
		} else if isValidMove(maze, previous, current, nextPoints.South, grid.South) {
			previous = current
			current = nextPoints.South
		} else if isValidMove(maze, previous, current, nextPoints.West, grid.West) {
			previous = current
			current = nextPoints.West
		} else {
			break
		}

		steps = append(steps, current)
	}

	return steps

}

func numIntersections(row []string, rowIdx int, pipePositionsInRow []grid.Point, forward bool) int {
	pipes := utils.Filter(pipePositionsInRow, func(p grid.Point) bool {
		if forward {
			return p.X > rowIdx && slices.Contains([]string{"F", "7", "L", "J", "|"}, row[p.X])
		} else {
			return p.X < rowIdx && slices.Contains([]string{"F", "7", "L", "J", "|"}, row[p.X])
		}
	})

	sort.Slice(pipes, func(i, j int) bool {
		return pipes[i].X < pipes[j].X
	})

	intersectChars := []string{}
	for _, p := range pipes {
		intersectChars = append(intersectChars, row[p.X])
	}

	intersectCount := 0
	for i := 0; i < len(intersectChars); i++ {
		if intersectChars[i] == "|" {
			intersectCount += 1
		}
		if i > 0 {
			if intersectChars[i-1] == "L" && intersectChars[i] == "7" {
				intersectCount += 1
			}
			if intersectChars[i-1] == "F" && intersectChars[i] == "J" {
				intersectCount += 1
			}
		}
	}

	return intersectCount
}

func Part1(input string) int {
	maze := grid.Parse(input)

	start := findStartPosition(maze)

	steps := getSteps(maze, start)

	return len(steps) / 2
}

func Part2(input string) int {
	maze := grid.Parse(input)

	start := findStartPosition(maze)

	steps := getSteps(maze, start)

	var pointsInsideTheLoop []grid.Point
	for y, row := range maze {

		pipePositionsInRow := utils.Filter(steps, func(p grid.Point) bool {
			return p.Y == y
		})

		if len(pipePositionsInRow) == 0 {
			// No pipes in this row - skip
			continue
		}

		for x := range row {

			if (slices.Contains(steps, grid.Point{X: x, Y: y})) {
				// This position is a pipe - ignore it
				continue
			}

			// Number of pipe intersections looking backwards
			backwardIntersectCount := numIntersections(row, x, pipePositionsInRow, false)

			// Number of pipe intersections looking forwards
			forwardIntersectCount := numIntersections(row, x, pipePositionsInRow, true)

			// If the number of intersections (for backwards and forwards) is odd, then the point is inside the loop
			if math.Mod(float64(backwardIntersectCount), 2) == 1 && math.Mod(float64(forwardIntersectCount), 2) == 1 {
				pointsInsideTheLoop = append(pointsInsideTheLoop, grid.Point{X: x, Y: y})
			}

		}
	}

	return len(pointsInsideTheLoop)
}
