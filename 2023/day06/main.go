package main

import (
	"fmt"
	"os"
	"strconv"
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

type race struct {
	time     int
	distance int
}

func Part1(input string) int {

	lines := strings.Split(input, "\n")

	times := parsing.ReadDigits(strings.Split(lines[0], ": ")[1])
	distances := parsing.ReadDigits(strings.Split(lines[1], ": ")[1])

	var races []race
	for i := 0; i < len(distances); i++ {
		races = append(races, race{time: times[i], distance: distances[i]})
	}

	var totalWaysToWin []int
	for _, race := range races {
		numWaysToWin := 0
		for i := 1; i < race.time; i++ {
			dist := i * (race.time - i)
			if dist > race.distance {
				numWaysToWin += 1
			}
		}
		totalWaysToWin = append(totalWaysToWin, numWaysToWin)
	}

	return utils.Product(totalWaysToWin)
}

func Part2(input string) int {
	lines := strings.Split(input, "\n")

	timeString := strings.ReplaceAll(strings.Split(lines[0], ": ")[1], " ", "")
	distanceString := strings.ReplaceAll(strings.Split(lines[1], ": ")[1], " ", "")

	time, _ := strconv.Atoi(timeString)
	distance, _ := strconv.Atoi(distanceString)

	race := race{time: time, distance: distance}

	numWaysToWin := 0
	for i := 1; i < race.time; i++ {
		dist := i * (race.time - i)
		if dist > race.distance {
			numWaysToWin += 1
		}
	}

	return numWaysToWin
}
