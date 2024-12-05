package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

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

func parse(input string) ([][]int, [][]int) {
	raw := strings.Split(strings.TrimSpace(input), "\n\n")

	var rules [][]int
	for _, rule := range parsing.ReadLines(raw[0]) {
		parts := strings.Split(rule, "|")
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		rules = append(rules, []int{a, b})
	}

	var allUpdates [][]int
	for _, update := range parsing.ReadLines(raw[1]) {
		var updates []int
		nums := strings.Split(update, ",")
		for _, num := range nums {
			n, _ := strconv.Atoi(num)
			updates = append(updates, n)
		}
		allUpdates = append(allUpdates, updates)
	}

	return rules, allUpdates
}

func sortUpdate(update []int, rules [][]int) []int {
	sorted := make([]int, len(update))
	copy(sorted, update)

	sort.Slice(sorted, func(i, j int) bool {
		for _, rule := range rules {
			if sorted[i] == rule[0] && sorted[j] == rule[1] {
				if i > j {
					return false
				}
			}
		}
		return true
	})

	return sorted
}

func isSorted(update []int, rules [][]int) bool {
	for _, rule := range rules {
		idx1 := slices.Index(update, rule[0])
		idx2 := slices.Index(update, rule[1])
		if (idx1 > -1 && idx2 > -1) && (idx1 > idx2) {
			return false
		}
	}
	return true
}

func Part1(input string) int {
	rules, updates := parse(input)

	var middleNumbers []int
	for _, update := range updates {
		if isSorted(update, rules) {
			middleNumbers = append(middleNumbers, update[len(update)/2])
		}
	}

	return utils.Sum(middleNumbers)

}

func Part2(input string) int {
	rules, updates := parse(input)

	middleNumbers := make([]int, 0)
	for _, update := range updates {
		if !isSorted(update, rules) {
			sorted := sortUpdate(update, rules)
			middleNumbers = append(middleNumbers, sorted[len(sorted)/2])
		}
	}

	return utils.Sum(middleNumbers)
}
