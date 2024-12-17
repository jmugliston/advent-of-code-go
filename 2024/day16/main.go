package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmugliston/aoc/grid"
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

type QueueItem struct {
	Position grid.PointWithDirection
	Path     []grid.PointWithDirection
	Score    int
}

// Calculate the smallest number of clockwise/anti-clockwise turns to reach the target direction
func shortestTurn(from grid.Direction, to grid.Direction) int {
	turns := map[grid.Direction]int{
		grid.North: 0,
		grid.East:  1,
		grid.South: 2,
		grid.West:  3,
	}

	diff := turns[to] - turns[from]

	if diff < 0 {
		diff += 4
	}

	if diff > 2 {
		return 4 - diff
	}

	return diff
}

func RunAlgortihm(maze grid.StringGrid, start grid.PointWithDirection, end grid.Point) (int, int) {
	// Map points visited to the number of steps taken to get there
	visited := make(map[grid.PointWithDirection]int)

	best := QueueItem{Score: 999999}
	pathMap := make(map[int][][]grid.PointWithDirection)

	queue := []QueueItem{{Position: start, Path: []grid.PointWithDirection{start}, Score: 0}}
	for len(queue) > 0 {
		current := queue[0]

		queue = queue[1:]

		if current.Position.X == end.X && current.Position.Y == end.Y {
			if current.Score <= best.Score {
				best = current
				pathMap[current.Score] = append(pathMap[current.Score], current.Path)
			}
			continue
		}

		if _, ok := visited[current.Position]; ok {
			if visited[current.Position] < current.Score {
				continue
			}
		}

		visited[current.Position] = current.Score

		// Check if we can move in any direction
		for _, direction := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {
			newPosition := grid.GetNextPointInDirection(grid.PointWithDirection{X: current.Position.X, Y: current.Position.Y, Direction: direction})

			if maze.GetPoint(newPosition) == "#" {
				// Hit a wall
				continue
			}

			newPositionWithDirection := grid.PointWithDirection{X: newPosition.X, Y: newPosition.Y, Direction: direction}

			// Make sure newPosition is not in the path we took
			if len(current.Path) > 1 {
				for _, step := range current.Path {
					if step == newPositionWithDirection {
						continue
					}
				}
			}

			newPath := make([]grid.PointWithDirection, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, newPositionWithDirection)

			turns := shortestTurn(current.Position.Direction, direction)

			queue = append(queue, QueueItem{Position: newPositionWithDirection, Path: newPath, Score: current.Score + 1 + (turns * 1000)})
		}
	}

	uniquePositions := make(map[grid.Point]bool)
	for _, path := range pathMap[best.Score] {
		for _, position := range path {
			uniquePositions[grid.Point{X: position.X, Y: position.Y}] = true
		}
	}

	return best.Score, len(uniquePositions)
}

func Part1(input string) int {
	maze := grid.Parse(input)

	startPoint := maze.Find("S")
	start := grid.PointWithDirection{X: startPoint.X, Y: startPoint.Y, Direction: grid.East}
	end := maze.Find("E")

	score, _ := RunAlgortihm(maze, start, end)

	return score
}

func Part2(input string) int {
	maze := grid.Parse(input)

	startPoint := maze.Find("S")
	start := grid.PointWithDirection{X: startPoint.X, Y: startPoint.Y, Direction: grid.East}
	end := maze.Find("E")

	_, seats := RunAlgortihm(maze, start, end)

	return seats
}