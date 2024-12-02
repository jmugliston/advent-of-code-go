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

	expected := 54

	result := Part1(string(input))

	if result != expected {
		t.Errorf("Expected %v, got %d", expected, result)
	}

}