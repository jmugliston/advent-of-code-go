package main

import (
	"os"
	"testing"
)

func TestPart1(t *testing.T) {

	input, err := os.ReadFile("./input/example.txt")

	if err != nil {
		panic("Couldn't find the example file!")
	}

	expected := 11

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

	expected := 31

	result := Part2(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}
