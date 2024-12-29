package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/parsing"
	"gonum.org/v1/gonum/stat/combin"
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
	name    string
	input1  *Wire
	input2  *Wire
	output  *Wire
	operand string
	hasRun  bool
}

func (g Gate) String() string {
	return fmt.Sprintf(
		"%s %s %s -> %s",
		g.input1.name,
		g.operand,
		g.input2.name,
		g.output.name,
	)
}

func parseInput(input string) ([]*Gate, string, string) {
	sections := strings.Split(input, "\n\n")

	initialInputs := parsing.ReadLines(sections[0])

	inputA := ""
	inputB := ""
	for _, input := range initialInputs {
		split := strings.Split(input, ": ")
		if strings.HasPrefix(split[0], "x") {
			inputA = inputA + split[1]
		} else {
			inputB = inputB + split[1]
		}
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
			wires[input1Name] = &Wire{name: input1Name, value: 0, activated: false}
		}

		if _, ok := wires[input2Name]; !ok {
			wires[input2Name] = &Wire{name: input2Name, value: 0, activated: false}
		}

		if _, ok := wires[outputName]; !ok {
			wires[outputName] = &Wire{name: outputName, value: 0, activated: false}
		}

		gates = append(gates, &Gate{
			name:    wires[input1Name].name + operand + wires[input2Name].name,
			input1:  wires[input1Name],
			input2:  wires[input2Name],
			output:  wires[outputName],
			operand: operand,
			hasRun:  false,
		})
	}

	return gates, inputA, inputB
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

func setInitialInputs(gates []*Gate, inputA string, inputB string) {
	wireMap := make(map[string]*Wire)

	for _, gate := range gates {
		wireMap[gate.input1.name] = gate.input1
		wireMap[gate.input2.name] = gate.input2
		wireMap[gate.output.name] = gate.output
	}

	for i := 0; i < len(inputA); i++ {
		xValue, _ := strconv.Atoi(string(inputA[i]))
		yValue, _ := strconv.Atoi(string(inputB[i]))

		index := fmt.Sprintf("%02d", i)
		wireMap["x"+index].value = xValue
		wireMap["x"+index].activated = true
		wireMap["y"+index].value = yValue
		wireMap["y"+index].activated = true
	}
}

func getGatesWithInputWire(gates []*Gate, wireName string) []*Gate {
	gatesWithInput := make([]*Gate, 0)
	for _, gate := range gates {
		if gate.input1.name == wireName || gate.input2.name == wireName {
			gatesWithInput = append(gatesWithInput, gate)
		}
	}
	return gatesWithInput
}

func getNumInputWires(gates []*Gate, gateToCheck *Gate) int {
	queue := make([]*Gate, 0)
	queue = append(queue, gateToCheck)

	dependentWires := make(map[*Wire]bool)
	for len(queue) > 0 {
		nextGate := queue[0]
		queue = queue[1:]
		for _, g := range gates {
			if g.output.name == nextGate.input1.name || g.output.name == nextGate.input2.name {
				if strings.HasPrefix(g.input1.name, "x") || strings.HasPrefix(g.input1.name, "y") {
					dependentWires[g.input1] = true
				}
				if strings.HasPrefix(g.input2.name, "x") || strings.HasPrefix(g.input2.name, "y") {
					dependentWires[g.input2] = true
				}
				if _, ok := dependentWires[g.output]; ok {
					continue
				}
				dependentWires[g.output] = true
				queue = append(queue, g)
			}
		}
	}

	return len(dependentWires)
}

func allGatesHaveRun(gates []*Gate) bool {
	for _, gate := range gates {
		if !gate.hasRun {
			return false
		}
	}
	return true
}

func copyGates(gates []*Gate) []*Gate {
	newGates := make([]*Gate, 0)
	newWireMap := make(map[string]*Wire)
	for _, gate := range gates {
		if _, ok := newWireMap[gate.input1.name]; !ok {
			newWireMap[gate.input1.name] = &Wire{name: gate.input1.name, value: gate.input1.value, activated: gate.input1.activated}
		}
		if _, ok := newWireMap[gate.input2.name]; !ok {
			newWireMap[gate.input2.name] = &Wire{name: gate.input2.name, value: gate.input2.value, activated: gate.input2.activated}
		}
		if _, ok := newWireMap[gate.output.name]; !ok {
			newWireMap[gate.output.name] = &Wire{name: gate.output.name, value: gate.output.value, activated: gate.output.activated}
		}
		newGates = append(newGates, &Gate{
			name:    gate.name,
			input1:  newWireMap[gate.input1.name],
			input2:  newWireMap[gate.input2.name],
			output:  newWireMap[gate.output.name],
			operand: gate.operand,
			hasRun:  gate.hasRun,
		})
	}
	return newGates
}

func runCircuit(origGates []*Gate, inputA string, inputB string) (int, string) {
	gates := copyGates(origGates)

	setInitialInputs(gates, inputA, inputB)

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

	i := 0
	runNext := gatesStart
	for {
		i += 1
		next := make([]*Gate, 0)

		if len(runNext) == 0 || i > len(gates) {
			break
		}

		for _, gate := range runNext {
			runGate(gate)
			next = append(next, getGatesWithInputWire(gates, gate.output.name)...)
		}

		if allGatesHaveRun(gatesWithOutput) {
			break
		}

		runNext = next
	}

	binaryString := ""
	for _, gate := range gatesWithOutput {
		binaryString = binaryString + strconv.Itoa(gate.output.value)
	}

	binaryValue, _ := strconv.ParseInt(binaryString, 2, 64)

	return int(binaryValue), binaryString
}

func swapGateOutputs(a *Gate, b *Gate) {
	temp := a.output
	a.output = b.output
	b.output = temp
}

func testCircuit(gates []*Gate) bool {
	testCases := [][]string{
		{"111111111111111111111111111111111111111111111", "111111111111111111111111111111111111111111111"},
		{"100000000000000000000000000000000000000000000", "111111111111111111111111111111111111111111111"},
	}

	testCasesExpected := []int{
		70368744177662,
		35184372088832,
	}

	for i, testCase := range testCases {
		result, _ := runCircuit(gates, testCase[0], testCase[1])
		if result != testCasesExpected[i] {
			return false
		}
	}

	return true
}

func Part1(input string) int {
	gates, inputA, inputB := parseInput(input)

	output, _ := runCircuit(gates, inputA, inputB)

	return output
}

func Part2(input string) string {

	gates, _, _ := parseInput(input)

	// Look for gates that are connected to z outputs that are not XOR gates
	badGates := make([]*Gate, 0)
	for _, gate := range gates {
		if strings.HasPrefix(gate.output.name, "z") {
			if gate.operand != "XOR" && gate.output.name != "z45" {
				badGates = append(badGates, gate)
			}
		}
	}

	swappedGates := []string{}

	// Swap the bad gates
	for _, badGate := range badGates {
		zNum, _ := strconv.Atoi(badGate.output.name[1:])
		for _, gate := range gates {
			// Look for another XOR gate that has the correct number of dependent input wires
			// (z number * 6 e.g. for z03 it should be 18)
			if gate.operand == "XOR" && !strings.HasPrefix(gate.output.name, "z") && getNumInputWires(gates, gate) == zNum*6 {
				swappedGates = append(swappedGates, badGate.output.name)
				swappedGates = append(swappedGates, gate.output.name)
				swapGateOutputs(badGate, gate)
			}
		}
	}

	// Brute force the last gate (slow but it works)...
	for _, combo := range combin.Combinations(len(gates), 2) {
		a := gates[combo[0]]
		b := gates[combo[1]]

		// Do not swap if it's a z output gate
		if strings.HasPrefix(a.output.name, "z") || strings.HasPrefix(b.output.name, "z") {
			continue
		}

		// Do not swap if it would create cycle
		if a.output.name == b.input1.name || a.output.name == b.input2.name || b.output.name == a.input1.name || b.output.name == a.input2.name {
			continue
		}

		swapGateOutputs(a, b)

		working := testCircuit(gates)

		if working {
			swappedGates = append(swappedGates, a.output.name)
			swappedGates = append(swappedGates, b.output.name)
			break
		}

		swapGateOutputs(a, b)
	}

	slices.Sort(swappedGates)

	return strings.Join(swappedGates, ",")
}
