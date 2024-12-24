package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/grid"
	"github.com/jmugliston/aoc/parsing"
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

type QueueItem struct {
	Point grid.Point
	Path  []grid.Point
	Keys  []string
}

func shortestPaths(pad grid.StringGrid, start string, end string) [][]string {
	startPoint := pad.Find(start)
	endPoint := pad.Find(end)

	queue := []QueueItem{{startPoint, []grid.Point{startPoint}, []string{}}}

	visited := map[grid.Point]bool{}

	best := math.MaxInt64
	paths := [][]string{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Point == endPoint {
			if len(current.Keys) <= best {
				best = len(current.Keys)
				paths = append(paths, append(current.Keys, "A"))
			}
			continue
		}

		visited[current.Point] = true

		for _, direction := range []grid.Direction{grid.North, grid.East, grid.South, grid.West} {
			next := current.Point.NextPoint(direction)

			if visited[next] {
				continue
			}

			if pad.GetPoint(next) == "." || !pad.IsPointInGrid(next) {
				continue
			}

			if len(current.Path)+1 > best {
				continue
			}

			nextPath := make([]grid.Point, len(current.Path))
			copy(nextPath, current.Path)
			nextPath = append(nextPath, next)

			nextKeys := make([]string, len(current.Keys))
			copy(nextKeys, current.Keys)

			var nextKey string
			switch direction {
			case grid.North:
				nextKey = "^"
			case grid.East:
				nextKey = ">"
			case grid.South:
				nextKey = "v"
			case grid.West:
				nextKey = "<"
			}

			nextKeys = append(nextKeys, nextKey)

			queue = append(queue, QueueItem{next, nextPath, nextKeys})
		}
	}

	finalPaths := [][]string{}
	for _, path := range paths {
		if len(path) <= best+1 {
			finalPaths = append(finalPaths, path)
		}
	}

	return finalPaths
}

func getSequences(sequence []string, level int, lookups map[string][][]string, cache map[string][][]string) [][]string {
	cacheKey := strings.Join(sequence, "") + strconv.Itoa(level)
	if cached, ok := cache[cacheKey]; ok {
		return cached
	}

	var nextSequences [][]string
	lastKey := "A"
	for _, key := range sequence {
		sPaths := lookups[lastKey+key]
		if len(nextSequences) == 0 {
			nextSequences = append(nextSequences, sPaths...)
		} else {
			var newSequences [][]string
			for _, seq := range nextSequences {
				for _, sPath := range sPaths {
					newSequences = append(newSequences, append(seq, sPath...))
				}
			}
			nextSequences = newSequences
		}
		lastKey = key
	}

	cache[cacheKey] = nextSequences

	if level < 2 {
		var next [][]string
		for _, seq := range nextSequences {
			next = append(next, getSequences(seq, level+1, lookups, cache)...)
		}
		cache[cacheKey] = next
		return next
	}

	return nextSequences
}

func Part1(input string) int {
	sequences := parsing.ReadLines(input)

	// 7 8 9
	// 4 5 6
	// 1 2 3
	// . 0 A
	numPad := grid.Parse("789\n456\n123\n.0A")

	// . ^ A
	// < v >
	dirPad := grid.Parse(".^A\n<v>")

	lookups := map[string][][]string{}

	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "A"}
	directions := []string{"^", "<", "v", ">", "A"}

	// Pre-compute all the shortest paths between numbers
	for _, numA := range numbers {
		for _, numB := range numbers {
			if numA == numB {
				lookups[numA+numB] = [][]string{{"A"}}
			}
			lookups[numA+numB] = shortestPaths(numPad, numA, numB)
		}
	}

	// Pre-compute all the shortest paths between directions
	for _, dirA := range directions {
		for _, dirB := range directions {
			if dirA == dirB {
				lookups[dirA+dirB] = [][]string{{"A"}}
			}
			lookups[dirA+dirB] = shortestPaths(dirPad, dirA, dirB)
		}
	}

	ans := 0
	for _, sequence := range sequences {
		cache := make(map[string][][]string, 0)

		allSequences := getSequences(strings.Split(sequence, ""), 0, lookups, cache)

		minLen := math.MaxInt64
		for _, sequence := range allSequences {
			if len(sequence) < minLen {
				minLen = len(sequence)
			}
		}

		num, _ := strconv.Atoi(sequence[:3])
		ans = ans + (minLen * num)

		fmt.Println(num, minLen, ans)
	}

	// 176452
	// not 176944
	return ans
}

func Part2(input string) int {
	return -1
}
