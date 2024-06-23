package main

import (
	"container/heap"
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

type CruciblePoint struct {
	grid.PointWithDirection
	Count int
}

func findPath(heatLossMap grid.NumberGrid, start grid.Point, end grid.Point, minCount int, maxCount int) int {

	visited := map[CruciblePoint]bool{}

	pq := make(PriorityQueue, 1)

	heap.Init(&pq)

	pq[0] = &QueueItem{
		PointWithDirection: grid.PointWithDirection{X: start.X, Y: start.Y, Direction: grid.North}, // Start with any direction that is not East/South
		Count:              minCount,
		Heatloss:           0,
		Path:               []grid.PointWithDirection{},
		Index:              0,
	}

	minHeatLoss := 9999999

	for pq.Len() > 0 {

		nextQueueItem := heap.Pop(&pq).(*QueueItem)

		visited[CruciblePoint{PointWithDirection: nextQueueItem.PointWithDirection, Count: nextQueueItem.Count}] = true

		currentStep := grid.PointWithDirection{X: nextQueueItem.X, Y: nextQueueItem.Y, Direction: nextQueueItem.Direction}
		currentCount := nextQueueItem.Count
		currentHeatLoss := nextQueueItem.Heatloss
		currentPath := nextQueueItem.Path

		if currentStep.X == end.X && currentStep.Y == end.Y {
			if currentCount < minCount {
				continue
			}
			minHeatLoss = currentHeatLoss
			break
		}

	DirectionLoop:
		for _, direction := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {

			nextStep := grid.GetNextPointInDirection(grid.PointWithDirection{X: currentStep.X, Y: currentStep.Y, Direction: direction})
			nextCount := currentCount

			if !heatLossMap.IsPointInGrid(nextStep) {
				continue
			}

			if direction == currentStep.Direction {
				nextCount += 1
				// Check for too many consecutive steps
				if nextCount > maxCount {
					continue
				}
			} else {
				// Check for too few consecutive steps
				if nextCount < minCount {
					continue
				}
				nextCount = 1
			}

			// Is already visited?
			if _, ok := visited[CruciblePoint{
				PointWithDirection: grid.PointWithDirection{X: nextStep.X, Y: nextStep.Y, Direction: direction},
				Count:              nextCount,
			}]; ok {
				continue
			}

			// Don't go back
			for _, step := range currentPath {
				if step.X == nextStep.X && step.Y == nextStep.Y {
					continue DirectionLoop
				}
			}

			// Already in the queue?
			for _, item := range pq {
				if item.X == nextStep.X && item.Y == nextStep.Y && item.Direction == direction && item.Count == nextCount {
					continue DirectionLoop
				}
			}

			nextPath := append([]grid.PointWithDirection(nil), currentPath...)
			nextPath = append(nextPath, grid.PointWithDirection{X: nextStep.X, Y: nextStep.Y, Direction: direction})

			heap.Push(&pq, &QueueItem{
				PointWithDirection: grid.PointWithDirection{X: nextStep.X, Y: nextStep.Y, Direction: direction},
				Count:              nextCount,
				Path:               nextPath,
				Heatloss:           currentHeatLoss + heatLossMap[nextStep.Y][nextStep.X],
			})
		}
	}

	return minHeatLoss
}

func Part1(input string) int {
	heatLossMap := grid.ParseNumbers(input)

	startPoint := grid.Point{X: 0, Y: 0}
	endPoint := grid.Point{X: len(heatLossMap[0]) - 1, Y: len(heatLossMap) - 1}

	heatLoss := findPath(heatLossMap, startPoint, endPoint, 0, 3)

	return heatLoss
}

func Part2(input string) int {
	heatLossMap := grid.ParseNumbers(input)

	startPoint := grid.Point{X: 0, Y: 0}
	endPoint := grid.Point{X: len(heatLossMap[0]) - 1, Y: len(heatLossMap) - 1}

	heatLoss := findPath(heatLossMap, startPoint, endPoint, 4, 10)

	return heatLoss
}
