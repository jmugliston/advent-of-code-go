package main

import (
	"os"
	"testing"
)

func TestPart1(t *testing.T) {

	input, err := os.ReadFile("./input/example2.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := 11687500

	result := Part1(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}

func TestPart2(t *testing.T) {

	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := 246006621493687

	result := Part2(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}
