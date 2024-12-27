package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}

func process(secret int) int {
	secret = prune(mix(secret, secret*64))

	secret = prune(mix(secret, secret/32))

	secret = prune(mix(secret, secret*2048))

	return secret
}

func Part1(input string) int {
	var nums []int

	for _, line := range parsing.ReadLines(input) {
		nums = append(nums, parsing.ReadNumbers(line)[0])
	}

	finalNums := []int{}
	for _, num := range nums {
		next := num
		for i := 0; i < 2000; i++ {
			next = process(next)
		}
		finalNums = append(finalNums, next)
	}

	result := 0
	for _, num := range finalNums {
		result += num
	}

	return result
}

func Part2(input string) int {
	var nums []int

	for _, line := range parsing.ReadLines(input) {
		nums = append(nums, parsing.ReadNumbers(line)[0])
	}

	prices := [][]int{}
	priceChanges := [][]int{}
	for _, num := range nums {
		next := num
		nextPrice := []int{next % 10}
		nextPriceChanges := []int{0}
		for i := 0; i < 2000; i++ {
			processed := process(next)
			nextPrice = append(nextPrice, processed%10)
			nextPriceChanges = append(nextPriceChanges, (processed%10)-nextPrice[i])
			next = processed
		}
		prices = append(prices, nextPrice)
		priceChanges = append(priceChanges, nextPriceChanges)
	}

	priceMap := map[string]int{}
	for num := range nums {
		seen := map[string]bool{}
		for i := 0; i < len(priceChanges[num])-4; i++ {
			window := priceChanges[num][i : i+4]
			key := fmt.Sprintf("%d,%d,%d,%d", window[0], window[1], window[2], window[3])
			if !seen[key] {
				seen[key] = true
				priceMap[key] += prices[num][i+3]
			}
		}
	}

	maximum := 0
	for _, price := range priceMap {
		if price > maximum {
			maximum = price
		}
	}

	return maximum
}
