package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/atheius/aoc/parsing"
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

type equation struct {
	result  int
	numbers []int
}

func parseInput(input string) []equation {
	lines := parsing.ReadLines(input)

	equations := make([]equation, 0)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		result, _ := strconv.Atoi(parts[0])
		numbers := parsing.ReadNumbers(parts[1])
		equations = append(equations, equation{
			result:  result,
			numbers: numbers,
		})
	}

	return equations
}

func isEquationValid(result int, numbers []int, useConcat bool) bool {
	if len(numbers) == 1 {
		return result == numbers[0]
	}

	nextNums := numbers[2:]

	if isEquationValid(result, append([]int{numbers[0] + numbers[1]}, nextNums...), useConcat) {
		return true
	}

	if isEquationValid(result, append([]int{numbers[0] * numbers[1]}, nextNums...), useConcat) {
		return true
	}

	if useConcat {
		concat, _ := strconv.Atoi(strconv.Itoa(numbers[0]) + strconv.Itoa(numbers[1]))
		return isEquationValid(result, append([]int{concat}, nextNums...), true)
	}

	return false
}

func Part1(input string) int {
	equations := parseInput(input)

	total := 0
	for _, equation := range equations {
		if isEquationValid(equation.result, equation.numbers, false) {
			total = total + equation.result
		}
	}

	return total
}

func Part2(input string) int {
	equations := parseInput(input)

	total := 0
	for _, equation := range equations {
		if isEquationValid(equation.result, equation.numbers, true) {
			total = total + equation.result
		}
	}

	return total
}
