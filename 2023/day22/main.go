package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

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

type Position struct {
	x int
	y int
	z int
}

type Brick struct {
	id          int
	positions   []Position
	bricksAbove []*Brick
	bricksBelow []*Brick
}

func parseBricks(lines []string) []*Brick {
	bricks := make([]*Brick, 0)
	for idx, line := range lines {
		var startX, startY, startZ int
		var endX, endY, endZ int

		parsed, _ := fmt.Sscanf(
			line,
			"%d,%d,%d~%d,%d,%d",
			&startX, &startY, &startZ,
			&endX, &endY, &endZ)

		if parsed != 6 {
			panic("Could not parse line")
		}

		brickPositions := make([]Position, 0)
		for x := startX; x <= endX; x++ {
			for y := startY; y <= endY; y++ {
				for z := startZ; z <= endZ; z++ {
					brickPositions = append(brickPositions, Position{x, y, z})
				}
			}
		}

		bricks = append(bricks, &Brick{
			id:        idx,
			positions: brickPositions,
		})
	}

	return bricks
}

// Create a map of all brick positions
func mapAllBrickPositions(bricks []*Brick) map[Position]*Brick {
	brickPositions := make(map[Position]*Brick)
	for _, brick := range bricks {
		for _, position := range brick.positions {
			brickPositions[position] = brick
		}
	}

	return brickPositions
}

// Check if a brick can fall by checking if the position below is empty
func canBrickFall(brick *Brick, positionMap map[Position]*Brick) bool {
	for _, position := range brick.positions {

		// Is position at the bottom?
		if (position.z - 1) <= 0 {
			return false
		}

		// Is position occupied by another brick
		if brickToCheck, ok := positionMap[Position{
			x: position.x,
			y: position.y,
			z: position.z - 1,
		}]; ok {
			if brickToCheck.id != brick.id {
				return false
			}
		}

	}
	return true
}

// Stabilise the bricks by moving them down until they can't fall any further
func stabiliseBricks(bricks []*Brick, brickPositions map[Position]*Brick) []*Brick {
	keepGoing := true

	for keepGoing {
		keepGoing = false
		for i := 0; i < len(bricks); i++ {
			canFall := canBrickFall(bricks[i], brickPositions)
			if canFall {
				keepGoing = true

				prevPositions := bricks[i].positions

				// Reset the brick positions
				bricks[i].positions = make([]Position, 0)

				for _, prevPosition := range prevPositions {

					// Remove the previous position
					delete(brickPositions, prevPosition)

					// Add the new position
					brickPositions[Position{
						x: prevPosition.x,
						y: prevPosition.y,
						z: prevPosition.z - 1,
					}] = bricks[i]

					// Remap the brick position (1 unit down)
					bricks[i].positions = append(bricks[i].positions, Position{
						x: prevPosition.x,
						y: prevPosition.y,
						z: prevPosition.z - 1,
					})
				}
			}
		}
	}

	return bricks
}

// Check each brick for other bricks directly above / below
// and add them as dependencies.
func mapBrickDependencies(bricks []*Brick) {
	for _, brick := range bricks {
		for _, brickPosition := range brick.positions {
			for _, otherBrick := range bricks {
				if brick.id == otherBrick.id {
					continue
				}

				if (slices.Contains(otherBrick.positions, Position{
					x: brickPosition.x,
					y: brickPosition.y,
					z: brickPosition.z + 1,
				})) {
					if !slices.Contains(brick.bricksAbove, otherBrick) {
						brick.bricksAbove = append(brick.bricksAbove, otherBrick)
					}
					if !slices.Contains(otherBrick.bricksBelow, brick) {
						otherBrick.bricksBelow = append(otherBrick.bricksBelow, brick)
					}
				}
			}
		}
	}
}

// Check if a brick can be safely removed (i.e. no bricks above it that require it for support)
func isBrickSafeToRemove(brick *Brick) bool {
	for _, brickAbove := range brick.bricksAbove {
		if len(brickAbove.bricksBelow) == 1 {
			return false
		}
	}
	return true
}

// Disintegrate a brick and all bricks above it
func disintegrateBricks(brick *Brick) int {
	disintegratedBricks := make([]*Brick, 0)

	disintegratedBricks = append(disintegratedBricks, brick)

	keepGoing := true
	for keepGoing {
		keepGoing = false
		for _, disintegratedBrick := range disintegratedBricks {
			// Check each of the bricks above the one we're removing
			for _, brickAbove := range disintegratedBrick.bricksAbove {
				// If the brick has not already been disintegrated
				if !slices.Contains(disintegratedBricks, brickAbove) {
					allBricksBelowDisintegrated := true
					for _, brickBelow := range brickAbove.bricksBelow {
						if !slices.Contains(disintegratedBricks, brickBelow) {
							allBricksBelowDisintegrated = false
						}
					}
					// If all bricks below have disintegrated
					if allBricksBelowDisintegrated {
						disintegratedBricks = append(disintegratedBricks, brickAbove)
						keepGoing = true
					}
				}
			}
		}
	}

	return len(disintegratedBricks)
}

func Part1(input string) int {

	lines := parsing.ReadLines(input)

	bricks := parseBricks(lines)

	brickPositions := mapAllBrickPositions(bricks)

	stableBricks := stabiliseBricks(bricks, brickPositions)

	mapBrickDependencies(stableBricks)

	safeToMoveCount := 0
	for _, brick := range stableBricks {
		if isBrickSafeToRemove(brick) {
			safeToMoveCount++
		}
	}

	return safeToMoveCount
}

func Part2(input string) int {
	lines := parsing.ReadLines(input)

	bricks := parseBricks(lines)

	brickPositions := mapAllBrickPositions(bricks)

	stableBricks := stabiliseBricks(bricks, brickPositions)

	mapBrickDependencies(stableBricks)

	totalBricksDisintegrated := 0
	for _, brick := range stableBricks {
		// Disintegrate the brick and all bricks above it (-1 to account for the brick itself)
		totalBricksDisintegrated += disintegrateBricks(brick) - 1
	}

	return totalBricksDisintegrated
}
