package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/atheius/aoc/parsing"
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
	id    int
	start Position
	end   Position
}

func parseBricks(lines []string) []Brick {
	bricks := make([]Brick, 0)
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

		bricks = append(bricks, Brick{
			id:    idx,
			start: Position{startX, startY, startZ},
			end:   Position{endX, endY, endZ},
		})
	}
	return bricks
}

func createPositionMap(bricks []Brick) map[Position]Brick {
	var positionMap = make(map[Position]Brick)

	for _, brick := range bricks {
		for x := brick.start.x; x <= brick.end.x; x++ {
			for y := brick.start.y; y <= brick.end.y; y++ {
				for z := brick.start.z; z <= brick.end.z; z++ {
					position := Position{x, y, z}
					positionMap[position] = brick
				}
			}
		}
	}

	return positionMap
}

func canBrickFall(positionMap *map[Position]Brick, brick *Brick) bool {
	for x := brick.start.x; x <= brick.end.x; x++ {
		for y := brick.start.y; y <= brick.end.y; y++ {
			for z := brick.start.z; z <= brick.end.z; z++ {
				if (z - 1) <= 0 {
					return false
				}
				if brickBelow, ok := (*positionMap)[Position{x, y, z - 1}]; ok {
					if brick.id != brickBelow.id {
						return false
					}
				}
			}
		}
	}

	return true
}

func removeBrickPositions(positionMap *map[Position]Brick, brick *Brick) {
	for x := brick.start.x; x <= brick.end.x; x++ {
		for y := brick.start.y; y <= brick.end.y; y++ {
			for z := brick.start.z; z <= brick.end.z; z++ {
				delete(*positionMap, Position{x, y, z})
			}
		}
	}
}

func addBrickPositions(positionMap *map[Position]Brick, brick *Brick) {
	for x := brick.start.x; x <= brick.end.x; x++ {
		for y := brick.start.y; y <= brick.end.y; y++ {
			for z := brick.start.z; z <= brick.end.z; z++ {
				(*positionMap)[Position{x, y, z}] = *brick
			}
		}
	}
}

func stabiliseBricks(brickMap map[int]Brick, positionMap map[Position]Brick) map[Position]Brick {

	keepGoing := true

	for keepGoing {
		keepGoing = false
		for i := 0; i < len(brickMap); i++ {
			brick := brickMap[i]
			canFall := canBrickFall(&positionMap, &brick)
			if canFall {
				keepGoing = true

				removeBrickPositions(&positionMap, &brick)

				// move brick down
				brick.start.z -= 1
				brick.end.z -= 1
				brickMap[i] = brick

				addBrickPositions(&positionMap, &brick)
			}
		}
	}

	return positionMap
}

func mapBrickDependencies(brickMap map[int]Brick, positionMap map[Position]Brick) (map[int]map[Brick]bool, map[int]map[Brick]bool) {
	above := make(map[int]map[Brick]bool)
	below := make(map[int]map[Brick]bool)

	for _, brick := range brickMap {
		above[brick.id] = make(map[Brick]bool)
		below[brick.id] = make(map[Brick]bool)
	}

	for _, brick := range brickMap {
		for x := brick.start.x; x <= brick.end.x; x++ {
			for y := brick.start.y; y <= brick.end.y; y++ {
				for z := brick.start.z; z <= brick.end.z; z++ {
					if otherBrick, ok := positionMap[Position{x, y, z + 1}]; ok {
						if otherBrick.id != brick.id {
							above[brick.id][otherBrick] = true
							below[otherBrick.id][brick] = true
						}
					}
				}
			}
		}
	}

	return above, below
}

func isBrickSafeToRemove(brick Brick, above map[int]map[Brick]bool, below map[int]map[Brick]bool) bool {
	for brickAbove := range above[brick.id] {
		if len(below[brickAbove.id]) == 1 {
			return false
		}
	}
	return true
}

func Part1(input string) int {

	lines := parsing.ReadLines(input)

	bricks := parseBricks(lines)

	brickMap := make(map[int]Brick)
	for _, brick := range bricks {
		brickMap[brick.id] = brick
	}

	positionMap := createPositionMap(bricks)

	stablePositionMap := stabiliseBricks(brickMap, positionMap)

	above, below := mapBrickDependencies(brickMap, stablePositionMap)

	safeToMoveCount := 0
	for _, brick := range brickMap {
		if isBrickSafeToRemove(brick, above, below) {
			safeToMoveCount++
		}
	}

	return safeToMoveCount
}

func Part2(input string) int {
	return -1
}
