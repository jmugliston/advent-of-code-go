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

func simulateGuard(guardMap grid.StringGrid, startPosition grid.Point, obstacle grid.Point) (bool, map[grid.Point]bool) {
	guardPosition := grid.PointWithDirection{X: startPosition.X, Y: startPosition.Y, Direction: grid.North}

	visitedPoints := make(map[grid.Point]bool)
	visitedPointsWithDirection := make(map[grid.PointWithDirection]int)

	visitedPoints[grid.Point{X: guardPosition.X, Y: guardPosition.Y}] = true
	visitedPointsWithDirection[guardPosition] = 1

	for {
		nextDirection := guardPosition.Direction
		nextPosition := grid.GetNextPointInDirection(guardPosition)

		if guardMap.GetPoint(nextPosition) == "#" || nextPosition == obstacle {
			// Guard hit something - turn right!
			nextDirection = nextDirection.TurnRight90()
			nextPosition = grid.Point{X: guardPosition.X, Y: guardPosition.Y}
		}

		guardPosition = grid.PointWithDirection{X: nextPosition.X, Y: nextPosition.Y, Direction: nextDirection}

		if !guardMap.IsPointInGrid(grid.Point{X: guardPosition.X, Y: guardPosition.Y}) {
			// Out of bounds - finished!
			break
		}

		visitedPoints[grid.Point{X: guardPosition.X, Y: guardPosition.Y}] = true
		visitedPointsWithDirection[guardPosition] += 1

		if visitedPointsWithDirection[guardPosition] > 1 {
			// In a loop!
			return true, visitedPoints
		}
	}

	return false, visitedPoints
}

func Part1(input string) int {
	guardMap := grid.Parse(input)

	startPosition := guardMap.Find("^")
	guardMap.SetPoint(startPosition, ".")

	_, visitedPoints := simulateGuard(guardMap, startPosition, grid.Point{X: -1, Y: -1})

	return len(visitedPoints)
}

func Part2(input string) int {
	guardMap := grid.Parse(input)

	startPosition := guardMap.Find("^")
	guardMap.SetPoint(startPosition, ".")

	_, potentialObstacles := simulateGuard(guardMap, startPosition, grid.Point{X: -1, Y: -1})

	loops := 0
	for location := range potentialObstacles {
		if location != startPosition {
			loop, _ := simulateGuard(guardMap, startPosition, location)
			if loop {
				loops++
			}
		}
	}

	return loops
}
