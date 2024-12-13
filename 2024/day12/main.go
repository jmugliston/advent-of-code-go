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

type Region struct {
	Name    string
	Tiles   []grid.Point
	Edges   int
	Corners int
}

func getTileEdges(tile grid.Point, regionMap grid.StringGrid) int {
	edges := 0

	points := grid.SurroundingPoints(tile)

	for _, point := range []grid.Point{points.North, points.East, points.South, points.West} {
		if regionMap.IsPointInGrid(point) {
			if regionMap[point.Y][point.X] != regionMap[tile.Y][tile.X] {
				edges++
			}
		} else {
			edges++
		}
	}

	return edges
}

func getTileCorners(tile grid.Point, regionMap grid.StringGrid) int {
	region := regionMap.GetPoint(tile)

	adjacentPoints := grid.SurroundingPoints(tile)

	n := regionMap.GetPoint(adjacentPoints.North)
	ne := regionMap.GetPoint(adjacentPoints.NorthEast)
	e := regionMap.GetPoint(adjacentPoints.East)
	se := regionMap.GetPoint(adjacentPoints.SouthEast)
	s := regionMap.GetPoint(adjacentPoints.South)
	sw := regionMap.GetPoint(adjacentPoints.SouthWest)
	w := regionMap.GetPoint(adjacentPoints.West)
	nw := regionMap.GetPoint(adjacentPoints.NorthWest)

	// 4 types of corner cases (can be inside and outside corners)

	// Top left
	// xx.
	// x..

	// Top right
	// .xx
	// ..x

	// Bottom left
	// x..
	// xx.

	// Bottom right
	// ..x
	// .xx

	matches := 0
	for _, point := range []grid.Point{adjacentPoints.North, adjacentPoints.East, adjacentPoints.South, adjacentPoints.West} {
		if regionMap.GetPoint(point) == region {
			matches++
		}
	}

	if matches == 1 {
		// Special case for end points (only 1 adjacent tile)
		return 2
	}

	corners := 0

	// Top left
	if e == region && s == region {
		if w != region && n != region {
			// Outside
			corners++
		}
		if se != region {
			// Inside
			corners++
		}
	}

	// Top right
	if w == region && s == region {
		if e != region && n != region {
			// Outside
			corners++
		}
		if sw != region {
			// Inside
			corners++
		}
	}

	// Bottom right
	if w == region && n == region {
		if s != region && e != region {
			// Outside
			corners++
		}
		if nw != region {
			// Inside
			corners++
		}
	}

	// Bottom left
	if e == region && n == region {
		if s != region && w != region {
			// Outside
			corners++
		}
		if ne != region {
			// Inside
			corners++
		}
	}

	return corners
}

func getRegion(startPoint grid.Point, regionMap grid.StringGrid) Region {
	visited := make(map[grid.Point]bool)

	region := Region{
		Name:  regionMap[startPoint.Y][startPoint.X],
		Tiles: []grid.Point{},
		Edges: 0,
	}

	queue := []grid.Point{startPoint}

	// Flood fill to get region
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true

		region.Tiles = append(region.Tiles, current)
		region.Edges += getTileEdges(current, regionMap)
		region.Corners += getTileCorners(current, regionMap)

		points := grid.SurroundingPoints(current)

		for _, point := range []grid.Point{points.North, points.East, points.South, points.West} {
			if regionMap.IsPointInGrid(point) && !visited[point] && regionMap[point.Y][point.X] == region.Name {
				queue = append(queue, point)
			}
		}
	}

	if len(region.Tiles) == 1 {
		region.Corners = 4
	}

	return region

}

func Part1(input string) int {
	regionMap := grid.Parse(input)

	total := 0

	visited := make(map[grid.Point]bool)
	for y, row := range regionMap {
		for x := range row {
			if visited[grid.Point{X: x, Y: y}] {
				continue
			}

			region := getRegion(grid.Point{X: x, Y: y}, regionMap)

			for _, point := range region.Tiles {
				visited[point] = true
			}

			total += region.Edges * len(region.Tiles)
		}
	}

	return total
}

func Part2(input string) int {
	regionMap := grid.Parse(input)

	total := 0

	visited := make(map[grid.Point]bool)
	for y, row := range regionMap {
		for x := range row {
			if visited[grid.Point{X: x, Y: y}] {
				continue
			}

			region := getRegion(grid.Point{X: x, Y: y}, regionMap)

			for _, point := range region.Tiles {
				visited[point] = true
			}

			total += region.Corners * len(region.Tiles)
		}
	}

	return total
}
