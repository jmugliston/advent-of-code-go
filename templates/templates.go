package templates

const MainTemplate = `package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	return -1
}

func Part2(input string) int {
	return -1
}

`

const TestTemplate = `package main

import (
	"os"
	"testing"
)

func TestPart1(t *testing.T) {

	input, err := os.ReadFile("./input/example.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := -1

	result := Part1(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}

func TestPart2(t *testing.T) {

	input, err := os.ReadFile("./input/example.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := -1

	result := Part2(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}
	
`
