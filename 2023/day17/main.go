package main

import (
	"container/heap"
	"fmt"
	"os"
	"strconv"

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

type CruciblePoint struct {
	X int
	Y int
	grid.Direction
	Count int
}

func getHeatLossAtPoint(heatLossMap *grid.Grid, point grid.Point) int {
	heatLoss, _ := strconv.Atoi((*heatLossMap)[point.Y][point.X])
	return heatLoss
}

func findPath(heatLossMap grid.Grid, start grid.Point, end grid.Point, minCount int, maxCount int) int {

	visited := map[CruciblePoint]bool{}

	pq := make(PriorityQueue, 1)

	heap.Init(&pq)

	pq[0] = &QueueItem{
		X:         start.X,
		Y:         start.Y,
		Direction: grid.North, // Start with any direction that is not East/South
		Count:     minCount,
		Heatloss:  0,
		Path:      []grid.PointWithDirection{},
		Index:     0,
	}

	minHeatLoss := 9999999

	for pq.Len() > 0 {

		nextQueueItem := heap.Pop(&pq).(*QueueItem)

		visited[CruciblePoint{X: nextQueueItem.X, Y: nextQueueItem.Y, Direction: nextQueueItem.Direction, Count: nextQueueItem.Count}] = true

		currentStep := grid.PointWithDirection{X: nextQueueItem.X, Y: nextQueueItem.Y, Direction: nextQueueItem.Direction}
		currentCount := nextQueueItem.Count
		currentHeatLoss := nextQueueItem.Heatloss
		currentPath := nextQueueItem.Path

		if currentStep.X == end.X && currentStep.Y == end.Y {
			if currentCount < minCount {
				continue
			}
			minHeatLoss = currentHeatLoss
			// minPath = currentPath
			break
		}

		for _, direction := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {

			nextStep := grid.GetNextPointInDirection(grid.PointWithDirection{X: currentStep.X, Y: currentStep.Y, Direction: direction})
			nextCount := currentCount

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

			// Is in grid?
			if !grid.IsPointInGrid(nextStep, heatLossMap) {
				continue
			}

			// Is already visited?
			if _, ok := visited[CruciblePoint{X: nextStep.X, Y: nextStep.Y, Direction: direction, Count: nextCount}]; ok {
				continue
			}

			// Don't go back
			movingBackwards := false
			for _, step := range currentPath {
				if step.X == nextStep.X && step.Y == nextStep.Y {
					movingBackwards = true
					break
				}
			}

			if movingBackwards {
				continue
			}

			// Calculate heat loss
			nextHeatLoss := currentHeatLoss + getHeatLossAtPoint(&heatLossMap, nextStep)

			// Already in the queue?
			alreadyInQueue := false
			for _, item := range pq {
				if item.X == nextStep.X && item.Y == nextStep.Y && item.Direction == direction && item.Count == nextCount {
					alreadyInQueue = true
					break
				}
			}

			if alreadyInQueue {
				continue
			}

			nextPath := make([]grid.PointWithDirection, len(currentPath))
			copy(nextPath, currentPath)
			nextPath = append(nextPath, grid.PointWithDirection{X: nextStep.X, Y: nextStep.Y, Direction: direction})

			// Add to queue
			heap.Push(&pq, &QueueItem{
				X:         nextStep.X,
				Y:         nextStep.Y,
				Direction: direction,
				Count:     nextCount,
				Path:      nextPath,
				Heatloss:  nextHeatLoss,
			})
		}
	}

	return minHeatLoss

}

func Part1(input string) int {
	heatLossMap := grid.Parse(input)

	startPoint := grid.Point{X: 0, Y: 0}
	endPoint := grid.Point{X: len(heatLossMap[0]) - 1, Y: len(heatLossMap) - 1}

	heatLoss := findPath(heatLossMap, startPoint, endPoint, 0, 3)

	return heatLoss
}

func Part2(input string) int {
	heatLossMap := grid.Parse(input)

	startPoint := grid.Point{X: 0, Y: 0}
	endPoint := grid.Point{X: len(heatLossMap[0]) - 1, Y: len(heatLossMap) - 1}

	heatLoss := findPath(heatLossMap, startPoint, endPoint, 4, 10)

	return heatLoss
}
