package main

import (
	"os"
	"slices"
	"testing"
)

func TestProgram(t *testing.T) {
	registers := make(map[string]int)

	registers["A"] = 0
	registers["B"] = 0
	registers["C"] = 9

	program := []int{2, 6}

	RunProgram(&registers, program)

	if registers["B"] != 1 {
		t.Errorf("Expected register B to be 1, got %d", registers["B"])
	}
}

func TestProgram2(t *testing.T) {
	registers := make(map[string]int)

	registers["A"] = 10
	registers["B"] = 0
	registers["C"] = 0

	program := []int{5, 0, 5, 1, 5, 4}

	output := RunProgram(&registers, program)

	if !slices.Equal(output, []int{0, 1, 2}) {
		t.Errorf("Expected 0,1,2 got %v", output)
	}
}

func TestProgram3(t *testing.T) {
	registers := make(map[string]int)

	registers["A"] = 2024
	registers["B"] = 0
	registers["C"] = 0

	program := []int{0, 1, 5, 4, 3, 0}

	output := RunProgram(&registers, program)

	if !slices.Equal([]int{4, 2, 5, 6, 7, 7, 7, 7, 3, 1, 0}, output) {
		t.Errorf("Expected 4,2,5,6,7,7,7,7,3,1,0 got %v", output)
	}

	if registers["A"] != 0 {
		t.Errorf("Expected register A to be 0, got %d", registers["A"])
	}
}

func TestProgram4(t *testing.T) {
	registers := make(map[string]int)

	registers["A"] = 0
	registers["B"] = 29
	registers["C"] = 0

	program := []int{1, 7}

	RunProgram(&registers, program)

	if registers["B"] != 26 {
		t.Errorf("Expected register B to be 26, got %d", registers["B"])
	}
}

func TestProgram5(t *testing.T) {
	registers := make(map[string]int)

	registers["A"] = 0
	registers["B"] = 2024
	registers["C"] = 43690

	program := []int{4, 0}

	RunProgram(&registers, program)

	if registers["B"] != 44354 {
		t.Errorf("Expected register B to be 44354, got %d", registers["B"])
	}
}

func TestPart1(t *testing.T) {

	input, err := os.ReadFile("./input/example.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := "4,6,3,5,6,3,5,2,1,0"

	result := Part1(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}

func TestPart2(t *testing.T) {

	input, err := os.ReadFile("./input/example2.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := "117440"

	result := Part2(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}

}
