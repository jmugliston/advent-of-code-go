package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jmugliston/aoc/grid"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")
var exampleFlag = flag.Bool("example", false, "Use the example instead of the puzzle input")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)

	inputFile := "input.txt"
	if *exampleFlag {
		inputFile = "example.txt"
	}
	path := filepath.Join(filepath.Dir(filename), "input", inputFile)

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

// Get a map of antenna names to their locations
func getAntennaLocations(antennaMap grid.StringGrid) map[string][]grid.Point {
	antennaLocations := make(map[string][]grid.Point)

	for i := 0; i < len(antennaMap); i++ {
		for j := 0; j < len(antennaMap[i]); j++ {
			if antennaMap[i][j] != "." {
				antennaLocations[antennaMap[i][j]] = append(antennaLocations[antennaMap[i][j]], grid.Point{X: j, Y: i})
			}
		}
	}

	return antennaLocations
}

// Find all possible pairs of points in a list
func getPairs(points []grid.Point) [][]grid.Point {
	pairMap := make(map[string]bool)
	var pairs [][]grid.Point

	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			pair := []grid.Point{points[i], points[j]}
			key := fmt.Sprintf("%v-%v", pair[0], pair[1])
			if !pairMap[key] {
				pairs = append(pairs, pair)
				pairMap[key] = true
			}
		}
	}

	return pairs
}

// Get a map of antenna pairs, where the key is the antenna name and the value is a list of pairs
func getAntennaPairsMap(antennaLocations map[string][]grid.Point) map[string][][]grid.Point {
	antennaPairs := make(map[string][][]grid.Point)

	for antenna, locations := range antennaLocations {
		if len(locations) < 2 {
			continue
		}
		pairs := getPairs(locations)
		antennaPairs[antenna] = pairs
	}

	return antennaPairs
}

// Find all the anti-nodes for a given set of antenna pairs
func getAntiNodes(antennaMap grid.StringGrid, antennaPairsMap map[string][][]grid.Point, includeAll bool) map[grid.Point]bool {
	antiNodes := map[grid.Point]bool{}

	for _, pairs := range antennaPairsMap {
		for _, pair := range pairs {
			dx := pair[1].X - pair[0].X
			dy := pair[1].Y - pair[0].Y

			if includeAll {
				// Include the antenna locations as anti-nodes
				antiNodes[pair[0]] = true
				antiNodes[pair[1]] = true
			}

			// Check diagonals at increasing distances, until both points are out of bounds
			i := 1
			for {
				nextDx := dx * i
				nextDy := dy * i

				// Find the next diagonal points at each end of the line
				diagonalPointA := grid.Point{X: pair[0].X - nextDx, Y: pair[0].Y - nextDy}
				diagonalPointB := grid.Point{X: pair[1].X + nextDx, Y: pair[1].Y + nextDy}

				aInbounds := antennaMap.IsPointInGrid(diagonalPointA)
				bInbounds := antennaMap.IsPointInGrid(diagonalPointB)

				if !aInbounds && !bInbounds {
					break
				}

				if aInbounds {
					antiNodes[diagonalPointA] = true
				}

				if bInbounds {
					antiNodes[diagonalPointB] = true
				}

				if !includeAll {
					// Only include the first diagonal anti-node for each pair
					break
				}

				i += 1
			}
		}
	}

	return antiNodes
}

func Part1(input string) int {
	antennaMap := grid.Parse(input)

	antennaLocations := getAntennaLocations(antennaMap)

	antennaPairsMap := getAntennaPairsMap(antennaLocations)

	antiNodes := getAntiNodes(antennaMap, antennaPairsMap, false)

	return len(antiNodes)
}

func Part2(input string) int {
	antennaMap := grid.Parse(input)

	antennaLocations := getAntennaLocations(antennaMap)

	antennaPairsMap := getAntennaPairsMap(antennaLocations)

	antiNodes := getAntiNodes(antennaMap, antennaPairsMap, true)

	return len(antiNodes)
}
