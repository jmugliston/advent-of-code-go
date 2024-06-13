package main

import (
	"fmt"
	"os"

	"github.com/atheius/aoc/grid"
)

func main() {
	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the input file!")
	}

	inputString := string(input)

	part1Answer := Part1(inputString)
	fmt.Println(part1Answer)

	part2Answer := Part2(inputString)
	fmt.Println(part2Answer)
}

type QueueItem struct {
	Point grid.Point
	Steps int
}

func FindNumberOfSteps(plotGrid grid.StringGrid, startPosition grid.Point, maxSteps int, includeOdd bool) map[grid.Point]int {
	stepMap := make(map[grid.Point]int)

	queue := []QueueItem{{Point: startPosition, Steps: 0}}

	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]

		if next.Steps > maxSteps {
			continue
		}

		if (next.Steps%2) == 0 || includeOdd {
			stepMap[grid.Point{X: next.Point.X, Y: next.Point.Y}] = next.Steps
		}

		nextPoints := grid.SurroundingPoints(next.Point)

		for _, point := range []grid.Point{nextPoints.North, nextPoints.East, nextPoints.South, nextPoints.West} {
			if !plotGrid.IsPointInGrid(point) {
				continue
			}
			if plotGrid[point.Y][point.X] == "#" {
				continue
			}
			if _, ok := stepMap[point]; ok {
				continue
			}
			existsInQueue := false
			for _, p := range queue {
				if p.Point == point {
					existsInQueue = true
					break
				}
			}
			if !existsInQueue {
				queue = append(queue, QueueItem{Point: point, Steps: next.Steps + 1})
			}
		}

	}

	return stepMap
}

func Part1(input string) int {
	plotGrid := grid.Parse(input)

	startPosition := plotGrid.Find("S")

	return len(FindNumberOfSteps(plotGrid, startPosition, 64, false))
}

func Part2(input string) int {
	plotGrid := grid.Parse(input)

	startPosition := plotGrid.Find("S")

	maxSteps := 26501365

	stepMap := FindNumberOfSteps(plotGrid, startPosition, maxSteps, true)

	// Due to the nature of the input (a repeating diamond shape), we can use a geometric solution to this problem.
	// This blog has a good explanation...
	// https://github.com/villuna/aoc23/wiki/A-Geometric-solution-to-advent-of-code-2023,-day-21

	halfGridLength := len(plotGrid) / 2

	evenCorners := 0
	oddCorners := 0

	evenFull := 0
	oddFull := 0

	for _, steps := range stepMap {
		if steps%2 == 0 {
			evenFull++
		} else {
			oddFull++
		}
		if steps > halfGridLength {
			if steps%2 == 0 {
				evenCorners++
			} else {
				oddCorners++
			}
		}
	}

	n := (maxSteps - halfGridLength) / len(plotGrid)

	total := (n+1)*(n+1)*oddFull + n*n*evenFull - (n+1)*oddCorners + n*evenCorners

	return total
}
