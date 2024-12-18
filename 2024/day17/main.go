package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
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

func parseInput(input string) (map[string]int, []int) {
	split := strings.Split(input, "\n\n")

	registers := make(map[string]int)
	program := make([]int, 0)

	registersRaw := split[0]
	for _, line := range strings.Split(registersRaw, "\n") {
		registerSplit := strings.Split(line, ": ")
		regsiterName := strings.Replace(registerSplit[0], "Register ", "", 1)
		registerValue, _ := strconv.Atoi(registerSplit[1])
		registers[regsiterName] = registerValue
	}

	programRaw := strings.Split(split[1], ": ")
	for _, nextRaw := range strings.Split(programRaw[1], ",") {
		nextValue, _ := strconv.Atoi(nextRaw)
		program = append(program, nextValue)
	}

	return registers, program
}

func getComboOperand(registers *map[string]int, operand int) int {
	if operand < 4 {
		return operand
	}

	if operand == 4 {
		return (*registers)["A"]
	}

	if operand == 5 {
		return (*registers)["B"]
	}

	if operand == 6 {
		return (*registers)["C"]
	}

	return -1
}

func adv(registers *map[string]int, operand int) int {
	val := (*registers)["A"] / (int(math.Pow(float64(2), float64(getComboOperand(registers, operand)))))
	(*registers)["A"] = val
	return val
}

func bxl(registers *map[string]int, operand int) int {
	val := (*registers)["B"] ^ operand
	(*registers)["B"] = val
	return val
}

func bst(registers *map[string]int, operand int) int {
	calc := getComboOperand(registers, operand) % 8

	base3Str := strconv.FormatInt(int64(calc), 3)

	for len(base3Str) < 3 {
		base3Str = "0" + base3Str
	}

	base3Str = base3Str[len(base3Str)-3:]

	base3Int, _ := strconv.ParseInt(base3Str, 3, 64)

	val := int(base3Int)

	(*registers)["B"] = val

	return val
}

func bxc(registers *map[string]int) int {
	value := (*registers)["B"] ^ (*registers)["C"]
	(*registers)["B"] = value
	return value
}

func jnz(instructionPointer *int, registers *map[string]int, operand int) bool {
	if (*registers)["A"] == 0 {
		return false
	}
	(*instructionPointer) = operand
	return true
}

func out(registers *map[string]int, operand int) int {
	return getComboOperand(registers, operand) % 8
}

func bdv(registers *map[string]int, operand int) int {
	val := (*registers)["A"] / (int(math.Pow(float64(2), float64(getComboOperand(registers, operand)))))
	(*registers)["B"] = val
	return val
}

func cdv(registers *map[string]int, operand int) int {
	calc := (int(math.Pow(float64(2), float64(getComboOperand(registers, operand)))))
	val := (*registers)["A"] / calc
	(*registers)["C"] = val
	return val
}

func RunProgram(registers *map[string]int, program []int) []int {
	var outputs []int

	instructionPointer := 0
	var opCode int
	var operand int

	for {
		if instructionPointer >= len(program)-1 {
			break
		}

		opCode = program[instructionPointer]
		operand = program[instructionPointer+1]

		jumped := false

		if opCode == 0 {
			adv(registers, operand)
		}

		if opCode == 1 {
			bxl(registers, operand)
		}

		if opCode == 2 {
			bst(registers, operand)
		}

		if opCode == 3 {
			jumped = jnz(&instructionPointer, registers, operand)
		}

		if opCode == 4 {
			bxc(registers)
		}

		if opCode == 5 {
			outVal := out(registers, operand)
			outputs = append(outputs, outVal)
		}

		if opCode == 6 {
			bdv(registers, operand)
		}

		if opCode == 7 {
			cdv(registers, operand)
		}

		if !jumped {
			instructionPointer += 2
		}
	}

	return outputs
}

func solve(n int, d int, program []int) int {
	res := []int{math.MaxInt64}

	if d == -1 {
		return n
	}

	// Work backwards through the digits to solve

	// Check each 3 bit value (numbers 0-7) because the program output uses mod 8
	for i := 0; i < 8; i++ {
		// Digit changes every 8^d iterations, where d is the index of the digit in the program
		nn := n + i*int(math.Pow(8, float64(d)))

		registers := make(map[string]int)
		registers["A"] = nn
		registers["B"] = 0
		registers["C"] = 0
		outputs := RunProgram(&registers, program)

		if len(outputs) != len(program) {
			continue
		}

		if outputs[d] == program[d] {
			// Partial match - continue to next digit
			res = append(res, solve(nn, d-1, program))
		}
	}

	return slices.Min(res)
}

func Part1(input string) string {
	registers, program := parseInput(input)

	outputs := RunProgram(&registers, program)

	var result []string

	for _, output := range outputs {
		result = append(result, strconv.Itoa(output))
	}

	return strings.Join(result, ",")
}

func Part2(input string) string {
	_, program := parseInput(input)

	ans := solve(0, len(program)-1, program)

	return strconv.Itoa(ans)
}
