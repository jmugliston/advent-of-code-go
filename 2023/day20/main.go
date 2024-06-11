package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atheius/aoc/parsing"
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

type module struct {
	name       string
	moduleType string
	state      int
	nextState  int
	inputState map[string]int
	outputs    []*module
}

func parseNodes(lines []string) map[string]*module {
	nodes := map[string]*module{}

	for _, line := range lines {
		parts := strings.Split(line, "->")

		input := strings.TrimSpace(parts[0])
		inputPrefix := input[0]
		inputName := strings.TrimLeft(input, "&%")

		if input == "broadcaster" {
			nodes["broadcaster"] = &module{
				name:       "broadcaster",
				moduleType: "b",
				state:      0,
				outputs:    []*module{},
			}
		} else {
			if _, ok := nodes[inputName]; ok {
				nodes[inputName].moduleType = string(inputPrefix)
			} else {
				nodes[inputName] = &module{
					name:       inputName,
					moduleType: string(inputPrefix),
					state:      0,
					inputState: map[string]int{},
					outputs:    []*module{}}
			}
		}

		outputs := strings.Split(strings.TrimSpace(parts[1]), ", ")

		for _, output := range outputs {
			if _, ok := nodes[output]; !ok {
				nodes[output] = &module{name: output, moduleType: "unknown", state: 0, inputState: map[string]int{}, outputs: []*module{}}
			}
			nodes[output].inputState[inputName] = 0
			nodes[inputName].outputs = append(nodes[inputName].outputs, nodes[output])
		}
	}

	return nodes
}

type pulseItem struct {
	src   *module
	dst   *module
	pulse int
}

type processResult struct {
	highPulses int
	lowPulses  int
	keyCycles  map[string]int
}

func runPulse(src *module, dst *module, pulse int) []pulseItem {
	nextPulses := []pulseItem{}

	if dst.moduleType == "b" {
		for _, output := range dst.outputs {
			nextPulses = append(nextPulses, pulseItem{src: dst, dst: output, pulse: pulse})
		}
	}

	if dst.moduleType == "%" {
		if pulse == 0 {
			for _, output := range dst.outputs {
				if dst.state == 0 {
					dst.nextState = 1
					nextPulses = append(nextPulses, pulseItem{src: dst, dst: output, pulse: 1})
				} else {
					dst.nextState = 0
					nextPulses = append(nextPulses, pulseItem{src: dst, dst: output, pulse: 0})
				}
			}
		}
	}

	if dst.moduleType == "&" {
		dst.inputState[src.name] = pulse

		allHigh := true
		for _, state := range dst.inputState {
			if state == 0 {
				allHigh = false
				break
			}
		}

		for _, output := range dst.outputs {
			if allHigh {
				nextPulses = append(nextPulses, pulseItem{src: dst, dst: output, pulse: 0})
			} else {
				nextPulses = append(nextPulses, pulseItem{src: dst, dst: output, pulse: 1})
			}
		}
	}

	return nextPulses
}

func runButtonCycle(nodes map[string]*module, keyNodeNames map[string]bool, cycleNum int) processResult {

	highPulses := 0
	lowPulses := 0
	keyCycles := map[string]int{}

	queue := []pulseItem{
		{src: nodes["button"], dst: nodes["broadcaster"], pulse: 0},
	}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if item.pulse == 1 {
			highPulses++
		} else {
			lowPulses++
			if val, ok := keyNodeNames[item.dst.name]; ok && val {
				keyCycles[item.dst.name] = cycleNum
			}
		}

		nextPulses := runPulse(item.src, item.dst, item.pulse)

		item.dst.state = item.dst.nextState

		queue = append(queue, nextPulses...)
	}

	return processResult{
		highPulses,
		lowPulses,
		keyCycles,
	}

}

func Part1(input string) int {

	lines := parsing.ReadLines(input)

	nodes := parseNodes(lines)

	nodes["button"] = &module{
		name:       "button",
		moduleType: "button",
		state:      0,
		outputs:    []*module{nodes["broadcaster"]},
	}

	highPulses := 0
	lowPulses := 0
	for i := 0; i < 1000; i++ {
		result := runButtonCycle(nodes, map[string]bool{}, 0)
		highPulses += result.highPulses
		lowPulses += result.lowPulses
	}

	return highPulses * lowPulses
}

func Part2(input string) int {
	lines := parsing.ReadLines(input)

	nodes := parseNodes(lines)

	// I built a graph using graphviz (see graph.png) to confirm that rx is only changed by
	// one conjunction (jm) which is changed by 4 other conjunctions. So we need to find when
	// the cycle time of when they each (independently) receive a low pulse, then find the
	// least common multiple of those cycles to know when jm receives a high pulse from
	// each 'key' node dh, sg, lm, db.
	keyNodeNames := map[string]bool{
		"dh": true,
		"sg": true,
		"lm": true,
		"db": true,
	}

	nodes["button"] = &module{
		name:       "button",
		moduleType: "button",
		state:      0,
		outputs:    []*module{nodes["broadcaster"]},
	}

	cycleMap := map[string][]int{}

	// Run enough cycles to find the cycle time of each key node
	for i := 0; i < 10000; i++ {
		result := runButtonCycle(nodes, keyNodeNames, i)

		for key, cycle := range result.keyCycles {
			cycleMap[key] = append(cycleMap[key], cycle)
		}
	}

	cycleLengths := []int{}
	for _, cycle := range cycleMap {
		cycleLengths = append(cycleLengths, cycle[1]-cycle[0])
	}

	return utils.LCM(cycleLengths)
}
