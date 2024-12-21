package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"

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
		fmt.Println(Part1(string(input), *exampleFlag))
	} else {
		fmt.Println(Part2(string(input), *exampleFlag))
	}
}

type QueueItem struct {
	Point grid.Point
	Time  int
	Path  []grid.PointWithDirection
}

// Find the path from start to end of maze (without any cheats)
func findPath(maze grid.StringGrid, start grid.Point, end grid.Point) []grid.PointWithDirection {
	startPath := []grid.PointWithDirection{{X: start.X, Y: start.Y, Direction: grid.South}}

	queue := []QueueItem{{start, 0, startPath}}

	final := QueueItem{
		Time: math.MaxInt64,
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Point == end {
			if final.Time == 0 || current.Time < final.Time {
				final = current
			}
			continue
		}

		for _, nextDirection := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {
			next := current.Point.NextPoint(nextDirection)
			if maze.GetPoint(next) == "." {

				if len(current.Path) >= 2 {
					if current.Path[len(current.Path)-2].X == next.X && current.Path[len(current.Path)-2].Y == next.Y {
						continue
					}
				}

				newPath := make([]grid.PointWithDirection, len(current.Path))
				copy(newPath, current.Path)
				newPath = append(newPath, grid.PointWithDirection{X: next.X, Y: next.Y, Direction: nextDirection})

				queue = append(queue, QueueItem{next, current.Time + 1, newPath})
			}
		}
	}

	return final.Path
}

// Check if we can cheat at the current point
func canCheat(current grid.Point, visited map[grid.Point]int, maxRange int) []int {
	possibleCheats := []int{}

	// Get all points within the maximum range
	pointsInRange := []grid.Point{}
	for x := current.X - maxRange; x <= current.X+maxRange; x++ {
		for y := current.Y - maxRange; y <= current.Y+maxRange; y++ {
			if x == current.X && y == current.Y {
				continue
			}
			if grid.ManhattenDistance(current, grid.Point{X: x, Y: y}) <= maxRange {
				if _, ok := visited[grid.Point{X: x, Y: y}]; ok {
					pointsInRange = append(pointsInRange, grid.Point{X: x, Y: y})
				}
			}
		}
	}

	for _, point := range pointsInRange {
		timeSaved := (visited[point] - visited[current]) - (grid.ManhattenDistance(point, current))
		if (timeSaved) > 0 {
			possibleCheats = append(possibleCheats, timeSaved)
		}
	}

	return possibleCheats
}

func getCheats(shortestPath []grid.PointWithDirection, maxCheat int) map[int]int {
	// Track visited points in the path and time taken to get there (without cheating)
	visited := make(map[grid.Point]int)
	for idx, step := range shortestPath {
		visited[grid.Point{X: step.X, Y: step.Y}] = idx
	}

	groupedByTimeSaved := make(map[int]int)

	// Iterate through points and group the cheats by time saved
	for _, point := range shortestPath {
		possibleCheats := canCheat(grid.Point{X: point.X, Y: point.Y}, visited, maxCheat)
		for _, cheatTimeSaved := range possibleCheats {
			groupedByTimeSaved[cheatTimeSaved]++
		}
	}

	return groupedByTimeSaved
}

func Part1(input string, example bool) int {
	maze := grid.Parse(input)

	start := maze.Find("S")
	end := maze.Find("E")

	maze.SetPoint(start, ".")
	maze.SetPoint(end, ".")

	path := findPath(maze, start, end)

	groupedByTimeSaved := getCheats(path, 2)

	timeSavedThreshold := 100

	if example {
		timeSavedThreshold = 0
	}

	result := 0
	for timeSaved, num := range groupedByTimeSaved {
		if (timeSaved) >= timeSavedThreshold {
			result += num
		}
	}

	return result
}

func Part2(input string, example bool) int {
	maze := grid.Parse(input)

	start := maze.Find("S")
	end := maze.Find("E")

	maze.SetPoint(start, ".")
	maze.SetPoint(end, ".")

	path := findPath(maze, start, end)

	groupedByTimeSaved := getCheats(path, 20)

	timeSavedThreshold := 100

	if example {
		timeSavedThreshold = 50
	}

	result := 0
	for timeSaved, num := range groupedByTimeSaved {
		if (timeSaved) >= timeSavedThreshold {
			result += num
		}
	}

	return result
}
