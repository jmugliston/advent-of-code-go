package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/atheius/aoc/parsing"
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

type race struct {
	time     int
	distance int
}

func Part1(input string) int {

	lines := strings.Split(input, "\n")

	times := parsing.ReadNumbers(strings.Split(lines[0], ": ")[1])
	distances := parsing.ReadNumbers(strings.Split(lines[1], ": ")[1])

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
