package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/atheius/aoc/utils"
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

func parseInput(input string) ([]string, map[string][]string) {

	lines := strings.Split(strings.TrimSpace(input), "\n")

	instructions := strings.Split(lines[0], "")

	nodes := make(map[string][]string)

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		split := strings.Split(line, " = ")
		node := split[0]
		leftAndRight := strings.NewReplacer("(", "", ",", "", ")", "").Replace(split[1])
		left := strings.Split(leftAndRight, " ")[0]
		right := strings.Split(leftAndRight, " ")[1]
		nodes[node] = []string{left, right}
	}

	return instructions, nodes

}

func getStepsToNode(start string, instructions []string, nodes map[string][]string, part2 bool) int {
	steps := 0

	currentNode := start
	for {
		nextInstruction := instructions[int(math.Mod(float64(steps), float64(len(instructions))))]

		if nextInstruction == "L" {
			currentNode = nodes[currentNode][0]
		} else {
			currentNode = nodes[currentNode][1]
		}

		steps += 1

		if part2 && strings.HasSuffix(currentNode, "Z") {
			break
		}

		if currentNode == "ZZZ" {
			break
		}
	}

	return steps
}

func Part1(input string) int {

	instructions, nodes := parseInput(input)

	steps := getStepsToNode("AAA", instructions, nodes, false)

	return steps
}

func Part2(input string) int {
	instructions, nodes := parseInput(input)

	var steps []int
	for node := range nodes {
		if strings.HasSuffix(node, "A") {
			steps = append(steps, getStepsToNode(node, instructions, nodes, true))
		}
	}

	return utils.LCM(steps)
}
