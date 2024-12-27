package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/parsing"
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

type Wire struct {
	name      string
	value     int
	activated bool
}

type Gate struct {
	input1  *Wire
	input2  *Wire
	output  *Wire
	operand string
	hasRun  bool
}

func (g Gate) String() string {
	return fmt.Sprintf(
		"%s (%v) %s %s (%v) -> %s (%v) - %v",
		g.input1.name,
		g.input1.value,
		g.operand,
		g.input2.name,
		g.input2.value,
		g.output.name,
		g.output.value,
		g.hasRun,
	)
}

func parseInput(input string) []*Gate {
	sections := strings.Split(input, "\n\n")

	initialInputs := parsing.ReadLines(sections[0])

	initialInputMap := make(map[string]int)
	for _, input := range initialInputs {
		split := strings.Split(input, ": ")
		value, _ := strconv.Atoi(split[1])
		initialInputMap[split[0]] = value
	}

	gatesRaw := parsing.ReadLines(sections[1])

	wires := make(map[string]*Wire)
	gates := make([]*Gate, 0)
	for _, gateRaw := range gatesRaw {
		split := strings.Split(gateRaw, " ")

		operand := split[1]

		input1Name := split[0]
		input2Name := split[2]
		outputName := split[4]

		if _, ok := wires[input1Name]; !ok {
			w := &Wire{name: input1Name, value: 0, activated: false}

			if v, ok := initialInputMap[input1Name]; ok {
				w.value = v
				w.activated = true
			}

			wires[input1Name] = w
		}

		if _, ok := wires[input2Name]; !ok {
			w := &Wire{name: input2Name, value: 0, activated: false}

			if v, ok := initialInputMap[input2Name]; ok {
				w.value = v
				w.activated = true
			}

			wires[input2Name] = w
		}

		if _, ok := wires[outputName]; !ok {
			wires[outputName] = &Wire{name: outputName, value: 0, activated: false}
		}

		gates = append(gates, &Gate{
			input1:  wires[input1Name],
			input2:  wires[input2Name],
			output:  wires[outputName],
			operand: operand,
			hasRun:  false,
		})
	}

	return gates
}

func runGate(gate *Gate) {
	if gate.hasRun {
		return
	}

	if !gate.input1.activated || !gate.input2.activated {
		return
	}

	switch gate.operand {
	case "AND":
		gate.output.value = gate.input1.value & gate.input2.value
	case "OR":
		gate.output.value = gate.input1.value | gate.input2.value
	case "XOR":
		gate.output.value = gate.input1.value ^ gate.input2.value
	}

	gate.output.activated = true
	gate.hasRun = true
}

func getGatesWithInputWire(gates []*Gate, wire *Wire) []*Gate {
	gatesWithInput := make([]*Gate, 0)
	for _, gate := range gates {
		if gate.input1 == wire || gate.input2 == wire {
			gatesWithInput = append(gatesWithInput, gate)
		}
	}
	return gatesWithInput
}

func allGatesHaveRun(gates []*Gate) bool {
	for _, gate := range gates {
		if !gate.hasRun {
			return false
		}
	}
	return true
}

func simulate(gates []*Gate) string {

	gatesStart := make([]*Gate, 0)
	gatesWithOutput := make([]*Gate, 0)

	for _, gate := range gates {
		if gate.input1.activated && gate.input2.activated {
			gatesStart = append(gatesStart, gate)
		}
		if strings.HasPrefix(gate.output.name, "z") {
			gatesWithOutput = append(gatesWithOutput, gate)
		}
	}

	sort.Slice(gatesWithOutput, func(i, j int) bool {
		return gatesWithOutput[i].output.name > gatesWithOutput[j].output.name
	})

	runNext := gatesStart
	for {
		next := make([]*Gate, 0)

		for _, gate := range runNext {
			runGate(gate)
			next = append(next, getGatesWithInputWire(gates, gate.output)...)
		}

		if allGatesHaveRun(gatesWithOutput) {
			break
		}

		runNext = next
	}

	output := ""
	for _, gate := range gatesWithOutput {
		output = output + strconv.Itoa(gate.output.value)
	}

	return output
}

func Part1(input string) int {
	gates := parseInput(input)

	output := simulate(gates)

	i, _ := strconv.ParseInt(output, 2, 64)

	return int(i)
}

func Part2(input string) int {
	return -1
}
