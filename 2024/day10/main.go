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

// Create a struct that can hold a point, and a path (slice of points)
type Path struct {
	Point grid.Point
	Path  []grid.Point
}

func Part1(input string) int {
	topographicMap := grid.ParseNumbers(input)

	// Find all the start positions
	startPositions := topographicMap.FindAll(0)

	total := 0
	for _, pos := range startPositions {

		visited := make(map[grid.Point]bool)
		possiblePaths := [][]grid.Point{}

		queue := []Path{
			{
				Point: pos,
				Path:  []grid.Point{pos},
			},
		}

		for len(queue) > 0 {
			currentPoint := queue[0].Point
			currentPath := queue[0].Path

			// Remove the first element from the queue
			queue = queue[1:]

			nextPositions := grid.SurroundingPoints(currentPoint)

			for _, nextPos := range []grid.Point{
				nextPositions.North,
				nextPositions.East,
				nextPositions.South,
				nextPositions.West,
			} {
				if !topographicMap.IsPointInGrid(nextPos) {
					continue
				}

				if topographicMap.GetPoint(nextPos) != topographicMap.GetPoint(currentPoint)+1 {
					continue
				}

				_, hasVisited := visited[nextPos]
				if hasVisited {
					continue
				}

				visited[nextPos] = true
				if topographicMap.GetPoint(nextPos) == 9 {
					possiblePaths = append(possiblePaths, append(currentPath, nextPos))
					continue
				}

				// Add next path to the queue
				queue = append(queue, Path{
					Point: nextPos,
					Path:  append(currentPath, nextPos),
				})
			}
		}

		total += len(possiblePaths)
	}

	return total
}

func Part2(input string) int {
	topographicMap := grid.ParseNumbers(input)

	// Find all the start positions
	startPositions := topographicMap.FindAll(0)

	total := 0
	for _, pos := range startPositions {

		possiblePaths := [][]grid.Point{}

		queue := []Path{
			{
				Point: pos,
				Path:  []grid.Point{pos},
			},
		}

		for len(queue) > 0 {
			currentPoint := queue[0].Point
			currentPath := queue[0].Path

			// Remove the first element from the queue
			queue = queue[1:]

			nextPositions := grid.SurroundingPoints(currentPoint)

			for _, nextPos := range []grid.Point{
				nextPositions.North,
				nextPositions.East,
				nextPositions.South,
				nextPositions.West,
			} {
				if !topographicMap.IsPointInGrid(nextPos) {
					continue
				}

				if topographicMap.GetPoint(nextPos) != topographicMap.GetPoint(currentPoint)+1 {
					continue
				}

				if topographicMap.GetPoint(nextPos) == 9 {
					possiblePaths = append(possiblePaths, append(currentPath, nextPos))
					continue
				}

				// Add next path to the queue
				queue = append(queue, Path{
					Point: nextPos,
					Path:  append(currentPath, nextPos),
				})
			}
		}

		total += len(possiblePaths)
	}

	return total
}
