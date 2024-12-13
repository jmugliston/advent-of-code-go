package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/jmugliston/aoc/grid"
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

func isJunction(hikeMap grid.StringGrid, point grid.Point) bool {
	pathCount := 0

	points := grid.SurroundingPoints(point)
	for _, point := range []grid.Point{
		points.North,
		points.East,
		points.South,
		points.West} {
		{
			if !hikeMap.IsPointInGrid(point) {
				continue
			}

			if hikeMap[point.Y][point.X] != "#" {
				pathCount++
			}

			if pathCount > 2 {
				return true
			}
		}

	}

	return false
}

func findJunctions(hikeMap grid.StringGrid) []grid.Point {
	junctions := make([]grid.Point, 0)

	for y, row := range hikeMap {
		for x, cell := range row {
			if cell != "#" {
				point := grid.Point{X: x, Y: y}
				if isJunction(hikeMap, point) {
					junctions = append(junctions, point)
				}

			}
		}
	}

	return junctions
}

type path struct {
	point grid.Point
	steps int
}

func getAdjacentJunctions(hikeMap grid.StringGrid, junctions []grid.Point, start grid.Point, slopes bool) []path {
	adjacentJunctions := make([]path, 0)

	type queueItem struct {
		point grid.Point
		steps int
		route []grid.Point
	}

	queue := []queueItem{{point: start, steps: 0, route: make([]grid.Point, 0)}}

	for len(queue) > 0 {

		item := queue[0]
		queue = queue[1:]

		isJunction := slices.Contains(junctions, item.point)
		isStart := item.point == start

		if !isStart && isJunction {
			adjacentJunctions = append(adjacentJunctions, path{point: item.point, steps: item.steps})
			continue
		}

		nextSteps := item.steps + 1
		nextRoute := append(item.route, item.point)

		points := grid.SurroundingPoints(item.point)

		for direction, point := range []grid.Point{
			points.North,
			points.East,
			points.South,
			points.West} {
			{
				if !hikeMap.IsPointInGrid(point) {
					continue
				}

				pointChar := hikeMap[point.Y][point.X]

				if pointChar == "#" {
					continue
				}

				if slices.Contains(nextRoute, point) {
					continue
				}

				if slopes && slices.Contains([]string{"<", ">", "v"}, pointChar) {
					if pointChar == "<" && direction == 1 {
						continue
					}
					if pointChar == ">" && direction == 3 {
						continue
					}
					if pointChar == "v" && direction == 0 {
						continue
					}
				}

				queue = append(queue, queueItem{point: point, steps: nextSteps, route: nextRoute})
			}
		}

	}

	return adjacentJunctions
}

func getLongestPath(junctionMap map[grid.Point][]path, start grid.Point, end grid.Point) int {

	type queueItem struct {
		point grid.Point
		steps int
		route []grid.Point
	}

	var longestPath queueItem

	queue := make([]queueItem, 0)

	queue = append(queue, queueItem{point: start, route: []grid.Point{start}, steps: 0})

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if item.point == end {
			if item.steps > longestPath.steps {
				longestPath = item
			}
			continue
		}

		for _, path := range junctionMap[item.point] {
			if !slices.Contains(item.route, path.point) {
				newRoute := make([]grid.Point, len(item.route))
				copy(newRoute, item.route)

				queue = append(
					queue,
					queueItem{
						point: path.point,
						route: append(newRoute, path.point),
						steps: item.steps + path.steps,
					})
			}
		}

	}

	return longestPath.steps
}

func createJunctionMap(hikeMap grid.StringGrid, junctions []grid.Point, slopes bool) map[grid.Point][]path {
	junctionMap := make(map[grid.Point][]path)

	for _, junction := range junctions {
		junctionMap[junction] = getAdjacentJunctions(hikeMap, junctions, junction, slopes)
	}

	return junctionMap
}

func Part1(input string) int {
	hikeMap := grid.Parse(input)

	start := grid.Point{X: 1, Y: 0}
	end := grid.Point{X: len(hikeMap[0]) - 2, Y: len(hikeMap) - 1}

	junctions := findJunctions(hikeMap)
	junctions = append(junctions, start, end)

	junctionMap := createJunctionMap(hikeMap, junctions, true)

	longestPath := getLongestPath(junctionMap, start, end)

	return longestPath
}

func Part2(input string) int {
	hikeMap := grid.Parse(input)

	start := grid.Point{X: 1, Y: 0}
	end := grid.Point{X: len(hikeMap[0]) - 2, Y: len(hikeMap) - 1}

	junctions := findJunctions(hikeMap)
	junctions = append(junctions, start, end)

	junctionMap := createJunctionMap(hikeMap, junctions, false)

	longestPath := getLongestPath(junctionMap, start, end)

	return longestPath
}
