package main

import (
	"fmt"
	"os"
	"sort"

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

func tiltNorth(rockMap grid.StringGrid) grid.StringGrid {

	height := len(rockMap)
	width := len(rockMap[0])

	tiltedMap := make(grid.StringGrid, height)

	for y := 0; y < height; y++ {
		if tiltedMap[y] == nil {
			tiltedMap[y] = make([]string, height)
		}
		for x := 0; x < width; x++ {
			if rockMap[y][x] == "O" {
				// Roll the rock as far North as possible
				for i := y; i >= 0; i-- {
					if (i == 0) || (tiltedMap[i-1][x] != ".") {
						tiltedMap[y][x] = "."
						tiltedMap[i][x] = "O"
						break
					}
				}
			} else {
				tiltedMap[y][x] = rockMap[y][x]
			}
		}
	}

	return tiltedMap
}

func calculateLoad(g grid.StringGrid) int {
	load := 0

	height := len(g)
	width := len(g[0])
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if g[y][x] == "O" {
				load = load + (height - y)
			}
		}
	}

	return load
}

func Part1(input string) int {

	rockMap := grid.Parse(input)

	tilted := tiltNorth(rockMap)

	totalLoad := calculateLoad(tilted)

	return totalLoad
}

func Part2(input string) int {

	rockMap := grid.Parse(input)

	loads := []int{}
	cycleMap := map[int][]int{}

	// Run for enough cycles for a pattern to emerge
	for i := 0; i < 200; i++ {
		// North
		rockMap = tiltNorth(rockMap)

		// West
		rockMap = rockMap.RotateClockwise()
		rockMap = tiltNorth(rockMap)

		// South
		rockMap = rockMap.RotateClockwise()
		rockMap = tiltNorth(rockMap)

		// East
		rockMap = rockMap.RotateClockwise()
		rockMap = tiltNorth(rockMap)

		// Rotate back to North
		rockMap = rockMap.RotateClockwise()

		currentLoad := calculateLoad(rockMap)
		loads = append(loads, currentLoad)

		cycleMap[currentLoad] = append(cycleMap[currentLoad], i)
	}

	keys := make([]int, len(cycleMap))

	i := 0
	for k := range cycleMap {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	firstLoopNumber := keys[0]

	nums := cycleMap[firstLoopNumber][:2]

	cycleLength := nums[1] - nums[0]

	firstTimeSeenInCycle := cycleMap[firstLoopNumber][0]

	offset := ((1_000_000_000 - firstTimeSeenInCycle) % cycleLength) - 1

	// Find the index of the firstTimeSeenInCycle in the loads array
	idx := -1
	for i, v := range loads {
		if v == firstLoopNumber {
			idx = i
			break
		}
	}

	return loads[idx+offset]
}
