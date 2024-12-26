package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
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

// BFS to find the shortest path between two keys on the keypad
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
				paths = append(paths, current.Keys)
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

// Build a map of all the shortest paths between keys
func buildKeyMap() map[string][][]string {

	// 7 8 9
	// 4 5 6
	// 1 2 3
	// . 0 A
	numPad := grid.Parse("789\n456\n123\n.0A")

	// . ^ A
	// < v >
	dirPad := grid.Parse(".^A\n<v>")

	keyMap := map[string][][]string{}

	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "A"}
	directions := []string{"^", "<", "v", ">", "A"}

	// Pre-compute all the shortest paths on the number pad
	for _, numA := range numbers {
		for _, numB := range numbers {
			if numA == numB {
				keyMap[numA+numB] = [][]string{{"A"}}
			}
			keyMap[numA+numB] = shortestPaths(numPad, numA, numB)
		}
	}

	// Pre-compute all the shortest paths on the directional pad
	for _, dirA := range directions {
		for _, dirB := range directions {
			if dirA == dirB {
				keyMap[dirA+dirB] = [][]string{{"A"}}
			}
			keyMap[dirA+dirB] = shortestPaths(dirPad, dirA, dirB)
		}
	}

	return keyMap
}

// Build all possible sequences of keys
func buildSequences(keys []string, index int, prevKey string, currPath []string, keyMap map[string][][]string) [][]string {
	if index == len(keys) {
		return [][]string{currPath}
	}

	var result [][]string
	for _, path := range keyMap[prevKey+keys[index]] {
		newPath := make([]string, len(currPath))

		copy(newPath, currPath)

		newPath = append(newPath, path...)
		newPath = append(newPath, "A")

		subResult := buildSequences(keys, index+1, keys[index], newPath, keyMap)
		result = append(result, subResult...)
	}

	return result
}

// Find the shortest sequence of keys
func shortestSequence(keys []string, depth int, cache map[string]int, keyMap map[string][][]string) int {
	if depth == 0 {
		return len(keys)
	}

	cacheKey := strings.Join(keys, "") + ":" + strconv.Itoa(depth)

	if cached, ok := cache[cacheKey]; ok {
		return cached
	}

	total := 0

	for _, subKey := range strings.SplitAfter(strings.Join(keys, ""), "A") {
		if subKey == "" {
			continue
		}

		sequences := buildSequences(strings.Split(subKey, ""), 0, "A", []string{}, keyMap)

		min := math.MaxInt64
		for _, sequence := range sequences {
			next := shortestSequence(sequence, depth-1, cache, keyMap)
			if next < min {
				min = next
			}
		}

		total = total + min
	}

	cache[cacheKey] = total

	return total
}

func Part1(input string) int {
	sequences := parsing.ReadLines(input)

	keyMap := buildKeyMap()

	levels := 2

	ans := 0
	for _, sequence := range sequences {
		keys := strings.Split(sequence, "")

		result := buildSequences(keys, 0, "A", []string{}, keyMap)
		shortestSeqs := []int{}

		for _, seq := range result {
			shortestSeqs = append(shortestSeqs, shortestSequence(seq, levels, make(map[string]int), keyMap))
		}

		minSeq := slices.Min(shortestSeqs)

		num, _ := strconv.Atoi(strings.Join(keys[:3], ""))

		ans = ans + (num * minSeq)
	}

	return ans
}

func Part2(input string) int {
	sequences := parsing.ReadLines(input)

	keyMap := buildKeyMap()

	levels := 25

	ans := 0
	for _, sequence := range sequences {
		keys := strings.Split(sequence, "")

		result := buildSequences(keys, 0, "A", []string{}, keyMap)
		shortestSeqs := []int{}

		for _, seq := range result {
			shortestSeqs = append(shortestSeqs, shortestSequence(seq, levels, make(map[string]int), keyMap))
		}

		minSeq := slices.Min(shortestSeqs)

		num, _ := strconv.Atoi(strings.Join(keys[:3], ""))

		ans = ans + (num * minSeq)
	}

	return ans
}
