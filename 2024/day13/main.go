package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/grid"
	"gonum.org/v1/gonum/mat"
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

type Game struct {
	A     grid.Point
	B     grid.Point
	Prize grid.Point
}

func parseInput(input string) []Game {
	games := make([]Game, 0)

	rawGames := strings.Split(strings.TrimSpace(input), "\n\n")

	for _, rawGame := range rawGames {
		lines := strings.Split(rawGame, "\n")
		var game Game
		for idx, line := range lines {
			split := strings.Split(line, ":")
			subSplit := strings.Split(split[1], ",")

			regex := regexp.MustCompile(`(X\+|Y\+|X=|Y=)(\d+)`)

			xStr := regex.FindStringSubmatch(subSplit[0])[2]
			yStr := regex.FindStringSubmatch(subSplit[1])[2]

			x, _ := strconv.Atoi(xStr)
			y, _ := strconv.Atoi(yStr)

			if idx == 0 {
				game.A = grid.Point{X: x, Y: y}
			} else if idx == 1 {
				game.B = grid.Point{X: x, Y: y}
			} else {
				game.Prize = grid.Point{X: x, Y: y}
			}
		}
		games = append(games, game)
	}

	return games
}

func solveLinearEquation(x1, x2, y1, y2, resX, resY float64) (float64, float64) {
	// Create a matrix A and a vector b
	A := mat.NewDense(2, 2, []float64{x1, x2, y1, y2})
	b := mat.NewVecDense(2, []float64{resX, resY})

	var x mat.VecDense
	if err := x.SolveVec(A, b); err != nil {
		panic("Could not solve the equations")
	}

	// Round to 3 decimal places to avoid precision errors
	finalX := math.Round(x.At(0, 0)*1000) / 1000
	finalY := math.Round(x.At(1, 0)*1000) / 1000

	return finalX, finalY
}

func Part1(input string) int {
	games := parseInput(input)

	total := 0
	for _, game := range games {
		x, y := solveLinearEquation(float64(game.A.X), float64(game.B.X), float64(game.A.Y), float64(game.B.Y), float64(game.Prize.X), float64(game.Prize.Y))

		// Both x and y must be integers for a valid solution
		if x == float64(int64(x)) && y == float64(int64(y)) {
			total += ((int(x) * 3) + int(y))
		}
	}

	return total
}

func Part2(input string) int {
	games := parseInput(input)

	total := 0
	for _, game := range games {
		game.Prize.X += 10000000000000
		game.Prize.Y += 10000000000000

		x, y := solveLinearEquation(float64(game.A.X), float64(game.B.X), float64(game.A.Y), float64(game.B.Y), float64(game.Prize.X), float64(game.Prize.Y))

		if x == float64(int64(x)) && y == float64(int64(y)) {
			total += ((int(x) * 3) + int(y))
		}
	}

	return total
}
