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

func generateCombinations(operators []string, n int) [][]string {
	if n == 0 {
		return [][]string{{}}
	}

	combinations := make([][]string, 0, int(math.Pow(float64(len(operators)), float64(n))))
	for _, combination := range generateCombinations(operators, n-1) {
		for _, operator := range operators {
			newCombination := append([]string{}, combination...)
			newCombination = append(newCombination, operator)
			combinations = append(combinations, newCombination)
		}
	}

	return combinations
}

func testEquation(e equation, operators []string) bool {
	n := len(e.numbers) - 1

	operatorCombinations := generateCombinations(operators, n)

	for _, combination := range operatorCombinations {

		result := e.numbers[0]
		for i := 1; i < len(e.numbers); i++ {
			operator := combination[i-1]
			nextNumber := e.numbers[i]

			switch operator {
			case "+":
				result = result + nextNumber
			case "*":
				result = result * nextNumber
			case "||":
				result, _ = strconv.Atoi(strconv.Itoa(result) + strconv.Itoa(nextNumber))
			}
		}

		if e.result == result {
			return true
		}

	}

	return false
}

func Part1(input string) int {
	operators := []string{"+", "*"}

	equations := parseInput(input)

	total := 0
	for _, equation := range equations {
		if testEquation(equation, operators) {
			total = total + equation.result
		}
	}

	return total
}

func Part2(input string) int {
	operators := []string{"+", "*", "||"}

	equations := parseInput(input)

	total := 0
	for _, equation := range equations {
		if testEquation(equation, operators) {
			total = total + equation.result
		}
	}

	return total
}
