package main

import (
	"fmt"
	"os"
	"slices"

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

func getNextSteps(floorMap grid.StringGrid, currentPointWithDirection grid.PointWithDirection) []grid.PointWithDirection {

	nextSteps := []grid.PointWithDirection{}

	nextpoint := grid.GetNextPointInDirection(currentPointWithDirection)

	if floorMap.IsPointInGrid(nextpoint) {
		// Empty space
		if floorMap[nextpoint.Y][nextpoint.X] == "." {
			nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: currentPointWithDirection.Direction})
		}

		// Mirror
		if floorMap[nextpoint.Y][nextpoint.X] == "/" {
			var turnedDirection grid.Direction

			switch currentPointWithDirection.Direction {
			case grid.North:
				turnedDirection = grid.East
			case grid.East:
				turnedDirection = grid.North
			case grid.South:
				turnedDirection = grid.West
			case grid.West:
				turnedDirection = grid.South
			}

			nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: turnedDirection})
		}

		// Mirror
		if floorMap[nextpoint.Y][nextpoint.X] == "\\" {
			var turnedDirection grid.Direction

			switch currentPointWithDirection.Direction {
			case grid.North:
				turnedDirection = grid.West
			case grid.East:
				turnedDirection = grid.South
			case grid.South:
				turnedDirection = grid.East
			case grid.West:
				turnedDirection = grid.North
			}

			nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: turnedDirection})
		}

		// Splitters
		if floorMap[nextpoint.Y][nextpoint.X] == "|" {
			// Did we come from East/West?
			if currentPointWithDirection.Direction == grid.East || currentPointWithDirection.Direction == grid.West {
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: grid.North})
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: grid.South})
			} else {
				// Otherwise carry on
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: currentPointWithDirection.Direction})
			}
		}

		if floorMap[nextpoint.Y][nextpoint.X] == "-" {
			// Did we come from North/South?
			if currentPointWithDirection.Direction == grid.North || currentPointWithDirection.Direction == grid.South {
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: grid.East})
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: grid.West})
			} else {
				// Otherwise carry on
				nextSteps = append(nextSteps, grid.PointWithDirection{X: nextpoint.X, Y: nextpoint.Y, Direction: currentPointWithDirection.Direction})
			}
		}
	}

	return nextSteps
}

func getEnergisedTiles(floorMap grid.StringGrid, start grid.PointWithDirection) int {

	positions := make(map[grid.PointWithDirection]bool)

	positionQueue := []grid.PointWithDirection{start}

	for {
		if len(positionQueue) == 0 {
			break
		}

		nextPosition := positionQueue[0]

		positionQueue = positionQueue[1:]

		nextSteps := getNextSteps(floorMap, nextPosition)

		for _, nextStep := range nextSteps {
			if _, ok := positions[nextStep]; !ok {
				positions[nextStep] = true
				positionQueue = append(positionQueue, nextStep)
			}
		}
	}

	tiles := make(map[grid.Point]bool)
	for position := range positions {
		tiles[grid.Point{X: position.X, Y: position.Y}] = true
	}

	return len(tiles)
}

func Part1(input string) int {

	floorMap := grid.Parse(input)

	startingPoint := grid.PointWithDirection{X: -1, Y: 0, Direction: grid.East}

	energisedTiles := getEnergisedTiles(floorMap, startingPoint)

	return energisedTiles
}

func Part2(input string) int {

	floorMap := grid.Parse(input)

	var energisedTilesList []int

	// Check each edge of the map
	for x := 0; x < len(floorMap[0]); x++ {
		// Top / Bottom
		energisedTilesList = append(energisedTilesList, getEnergisedTiles(floorMap, grid.PointWithDirection{X: x, Y: -1, Direction: grid.South}))
		energisedTilesList = append(energisedTilesList, getEnergisedTiles(floorMap, grid.PointWithDirection{X: x, Y: +1, Direction: grid.North}))
	}

	for y := 0; y < len(floorMap); y++ {
		// Left / Right
		energisedTilesList = append(energisedTilesList, getEnergisedTiles(floorMap, grid.PointWithDirection{X: -1, Y: y, Direction: grid.East}))
		energisedTilesList = append(energisedTilesList, getEnergisedTiles(floorMap, grid.PointWithDirection{X: len(floorMap[0]), Y: y, Direction: grid.West}))
	}

	return slices.Max(energisedTilesList)
}