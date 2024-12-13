package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

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

type Stone struct {
	value  int
	blinks int
}

var cache = make(map[Stone]int)

func runRules(stone Stone, maxBlinks int) int {
	if stone.blinks == maxBlinks {
		return 1
	}

	if val, ok := cache[stone]; ok {
		return val
	}

	var result int

	numString := strconv.Itoa(stone.value)

	if stone.value == 0 {
		// Rule #1
		result = runRules(Stone{value: 1, blinks: stone.blinks + 1}, maxBlinks)
	} else if len(numString)%2 == 0 {
		// Rule #2
		newStoneValue1, _ := strconv.Atoi(numString[:len(numString)/2])
		newStoneValue2, _ := strconv.Atoi(numString[len(numString)/2:])
		result = runRules(Stone{value: newStoneValue1, blinks: stone.blinks + 1}, maxBlinks) + runRules(Stone{value: newStoneValue2, blinks: stone.blinks + 1}, maxBlinks)
	} else {
		// Rule #3
		result = runRules(Stone{value: stone.value * 2024, blinks: stone.blinks + 1}, maxBlinks)
	}

	cache[stone] = result

	return result
}

func Part1(input string) int {
	stones := parsing.ReadNumbers(input)

	total := 0
	for i := 0; i < len(stones); i++ {
		total += runRules(Stone{value: stones[i], blinks: 0}, 25)
	}

	return total
}

func Part2(input string) int {
	stones := parsing.ReadNumbers(input)

	total := 0
	for i := 0; i < len(stones); i++ {
		total += runRules(Stone{value: stones[i], blinks: 0}, 75)
	}

	return total
}
