package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atheius/aoc/utils"
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
