package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
		fmt.Println(Part1(string(input), *exampleFlag))
	} else {
		fmt.Println(Part2(string(input), *exampleFlag))
	}
}

type QueueItem struct {
	Point grid.Point
	Path  []grid.PointWithDirection
	Steps int
}

func parseInput(input string) []grid.Point {
	points := make([]grid.Point, 0)

	lines := parsing.ReadLines(input)
	for _, line := range lines {
		split := strings.Split(line, ",")
		x, _ := strconv.Atoi(split[0])
		y, _ := strconv.Atoi(split[1])
		points = append(points, grid.Point{X: x, Y: y})
	}

	return points
}

func ShortestPath(memoryMap grid.StringGrid, start grid.Point, end grid.Point) []grid.PointWithDirection {
	queue := make([]QueueItem, 0)
	queue = append(queue, QueueItem{Point: start, Path: []grid.PointWithDirection{}, Steps: 0})

	best := QueueItem{Steps: math.MaxInt64}
	visited := make(map[grid.Point]int)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Point == end {
			if current.Steps < best.Steps || best.Steps == 0 {
				best = current
			}
			continue
		}

		if visited[current.Point] == 0 {
			visited[current.Point] = current.Steps
		} else if current.Steps < visited[current.Point] {
			visited[current.Point] = current.Steps
		} else {
			continue
		}

		for _, direction := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {

			nextPoint := current.Point.GetNextPointInDirection(direction)

			if memoryMap.GetPoint(nextPoint) != "." {
				continue
			}

			// Already visited this point?
			if steps, ok := visited[nextPoint]; ok {
				// Did we get here in less steps?
				if current.Steps+1 >= steps {
					continue
				}
			}

			// Check we're not going back on ourselves
			alreadyInPath := false
			for _, point := range current.Path {
				if point.X == nextPoint.X && point.Y == nextPoint.Y {
					alreadyInPath = true
					break
				}
			}

			if alreadyInPath {
				continue
			}

			newPath := make([]grid.PointWithDirection, len(current.Path))
			copy(newPath, current.Path)
			newPath = append(newPath, grid.PointWithDirection{X: nextPoint.X, Y: nextPoint.Y, Direction: direction})

			queue = append(queue, QueueItem{Point: nextPoint, Path: newPath, Steps: current.Steps + 1})
		}

	}

	return best.Path
}

func Part1(input string, example bool) int {
	points := parseInput(input)

	numBytes := 1024
	height := 71
	width := 71

	if example {
		numBytes = 12
		height = 7
		width = 7
	}

	memoryMap := grid.InitialiseStringGrid(height, width, ".")

	for idx, point := range points {
		if idx < numBytes {
			memoryMap.SetPoint(point, string("#"))
		}
	}

	path := ShortestPath(memoryMap, grid.Point{X: 0, Y: 0}, grid.Point{X: width - 1, Y: height - 1})

	return len(path)
}

func Part2(input string, example bool) string {
	points := parseInput(input)

	numBytes := 1024
	height := 71
	width := 71

	if example {
		numBytes = 12
		height = 7
		width = 7
	}

	memoryMap := grid.InitialiseStringGrid(height, width, ".")

	for idx, point := range points {
		if idx < numBytes {
			memoryMap.SetPoint(point, string("#"))
		}
	}

	var finalPixel grid.Point

	previousPath := ShortestPath(memoryMap, grid.Point{X: 0, Y: 0}, grid.Point{X: width - 1, Y: height - 1})

	for i := numBytes + 1; i < len(points); i++ {
		memoryMap.SetPoint(points[i], string("#"))

		// Check if next point is in previous path
		for _, point := range previousPath {
			if (points[i] == grid.Point{X: point.X, Y: point.Y}) {
				previousPath = ShortestPath(memoryMap, grid.Point{X: 0, Y: 0}, grid.Point{X: width - 1, Y: height - 1})
			}
		}

		if len(previousPath) == 0 {
			// Found the pixel that blocks the exit
			finalPixel = points[i]
			break
		}

	}

	return fmt.Sprintf("%d,%d", finalPixel.X, finalPixel.Y)
}
