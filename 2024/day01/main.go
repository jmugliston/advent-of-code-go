package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/jmugliston/aoc/parsing"
	"github.com/jmugliston/aoc/utils"
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

func Part1(input string) int {
	numbers := parsing.ReadNumbers(input)

	list_a := utils.EveryNthElement(numbers, 2)
	list_b := utils.EveryNthElement(numbers[1:], 2)

	sort.Ints(list_a)
	sort.Ints(list_b)

	total := 0
	for i := range list_a {
		total += utils.Abs(list_a[i] - list_b[i])
	}

	return total
}

func Part2(input string) int {
	numbers := parsing.ReadNumbers(input)

	list_a := utils.EveryNthElement(numbers, 2)

	freq_map_b := map[int]int{}
	for _, num := range utils.EveryNthElement(numbers[1:], 2) {
		freq_map_b[num]++
	}

	similarityScore := 0
	for _, num := range list_a {
		similarityScore += num * freq_map_b[num]
	}

	return similarityScore
}
