package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atheius/aoc/grid"
	"github.com/atheius/aoc/parsing"
)

func main() {
	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the input file!")
	}

	inputString := string(input)

	part1Answer := Part1(inputString)
	fmt.Println(part1Answer)

	part2Answer := Part2(inputString)
	fmt.Println(part2Answer)
}

type Instruction struct {
	grid.Direction
	Amount int
	Colour string
}

func getPerimiterPoints(instructions []Instruction) []grid.Point {
	var perimiterPoints []grid.Point

	currentPoint := grid.Point{X: 0, Y: 0}
	for _, i := range instructions {
		for j := 0; j < i.Amount; j++ {
			currentPoint = grid.GetNextPointInDirection(grid.PointWithDirection{X: currentPoint.X, Y: currentPoint.Y, Direction: i.Direction})
			perimiterPoints = append(perimiterPoints, currentPoint)
		}
	}

	return perimiterPoints
}

func getArea(perimiterPoints []grid.Point) int {
	area := grid.ShoelaceFormula(perimiterPoints)

	// Because the polygon is on a grid and the coordinates are the middle of squares,
	// we need to use Pick's Theorem to calculate the internal area...
	// A = i + b/2 - 1
	// i = A - b/2 + 1
	internalArea := area - (len(perimiterPoints) / 2) + 1

	return internalArea + len(perimiterPoints)
}

func Part1(input string) int {

	lines := parsing.ReadLines(input)

	instructions := make([]Instruction, 0)
	for _, line := range lines {
		split := strings.Split(line, " ")

		var direction grid.Direction

		switch split[0] {
		case "R":
			direction = grid.East
		case "D":
			direction = grid.South
		case "L":
			direction = grid.West
		case "U":
			direction = grid.North

		}

		amount, _ := strconv.Atoi(split[1])
		colour := strings.Trim(split[2], "()")

		instructions = append(instructions, Instruction{
			Direction: direction,
			Amount:    amount,
			Colour:    colour,
		})
	}

	perimiterPoints := getPerimiterPoints(instructions)

	return getArea(perimiterPoints)
}

func Part2(input string) int {
	lines := parsing.ReadLines(input)

	instructions := make([]Instruction, 0)
	for _, line := range lines {
		split := strings.Split(line, " ")

		hex := strings.Trim(split[2], "()#")

		lastChar := hex[len(hex)-1:]

		amount, _ := strconv.ParseInt(hex[:len(hex)-1], 16, 64)

		var direction grid.Direction
		switch lastChar {
		case "0":
			direction = grid.East
		case "1":
			direction = grid.South
		case "2":
			direction = grid.West
		case "3":
			direction = grid.North
		}

		instructions = append(instructions, Instruction{
			Direction: direction,
			Amount:    int(amount),
		})
	}

	perimiterPoints := getPerimiterPoints(instructions)

	return getArea(perimiterPoints)
}
