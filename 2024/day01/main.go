package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/atheius/aoc/parsing"
	"github.com/atheius/aoc/utils"
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

func Part1(input string) int {
	numbers := parsing.ReadDigits(input)

	list_a := []int{}
	list_b := []int{}

	for idx, num := range numbers {
		if idx%2 == 0 {
			list_a = append(list_a, num)
		} else {
			list_b = append(list_b, num)
		}
	}

	sort.Ints(list_a)
	sort.Ints(list_b)

	total := 0
	for i := range list_a {
		total += utils.Abs(list_a[i] - list_b[i])
	}

	return total
}

func Part2(input string) int {
	numbers := parsing.ReadDigits(input)

	list_a := []int{}
	freq_map_b := map[int]int{}

	for idx, num := range numbers {
		if idx%2 == 0 {
			list_a = append(list_a, num)
		} else {
			freq_map_b[num]++
		}
	}

	similarityScore := 0

	for _, num := range list_a {
		similarityScore += num * freq_map_b[num]
	}

	return similarityScore
}
