package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func HASHAlgortihm(input string) int {
	currentValue := 0
	for _, char := range input {
		currentValue = currentValue + int(char)
		currentValue = currentValue * 17
		currentValue = int(math.Mod(float64(currentValue), 256))
	}
	return currentValue
}

func Part1(input string) int {

	steps := strings.Split(strings.TrimSpace(input), ",")

	value := 0
	for _, step := range steps {
		value = value + HASHAlgortihm(step)
	}

	return value
}

func Part2(input string) int {
	steps := strings.Split(strings.TrimSpace(input), ",")

	type lens struct {
		label string
		value int
	}

	boxes := make(map[int][]lens, 250)

	for _, step := range steps {
		split := regexp.MustCompile("(=|-)").Split(step, -1)

		label := split[0]
		boxNumber := HASHAlgortihm(split[0])
		value, _ := strconv.Atoi(split[1])

		if strings.Contains(step, "=") {
			boxIdx := -1
			for idx, lens := range boxes[boxNumber] {
				if lens.label == label {
					boxIdx = idx
				}
			}

			if boxIdx != -1 {
				// Update the value of the lens in the box
				boxes[boxNumber][boxIdx].value = value
			} else {

				if len(boxes[boxNumber]) == 0 {
					boxes[boxNumber] = []lens{}
				}

				// Add the lens to the box
				boxes[boxNumber] = append(boxes[boxNumber], lens{label: label, value: value})
			}

		} else {
			for idx, lens := range boxes[boxNumber] {
				if lens.label == label {
					// Remove the lens from the box
					boxes[boxNumber] = append(boxes[boxNumber][:idx], boxes[boxNumber][idx+1:]...)
				}
			}
		}
	}

	focussingPower := 0
	for boxIdx, box := range boxes {
		for lensIdx, lens := range box {
			focussingPower = focussingPower + ((boxIdx + 1) * (lensIdx + 1) * lens.value)
		}
	}

	return focussingPower
}
