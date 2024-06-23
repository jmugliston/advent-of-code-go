package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atheius/aoc/parsing"
	"github.com/juliangruber/go-intersect"
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

func createCards(lines []string) []Card {

	cards := make([]Card, 0)

	for _, line := range lines {
		gameString := strings.Split(line, ": ")
		numbers := strings.Split(gameString[1], " | ")

		winningNumbers := parsing.ReadDigits(numbers[0])
		scratchcardNumbers := parsing.ReadDigits(numbers[1])
		points := len(intersect.Hash(winningNumbers, scratchcardNumbers))

		nextCard := Card{
			WinningNumbers:     winningNumbers,
			ScratchcardNumbers: scratchcardNumbers,
			Points:             points,
			Copies:             1,
		}

		cards = append(cards, nextCard)
	}

	return cards

}

func Part1(input string) int {

	lines := strings.Split(strings.TrimSpace(input), "\n")

	points := 0

	cards := createCards(lines)

	for _, card := range cards {
		if card.Points == 1 {
			points += 1
		} else if card.Points >= 1 {
			points += int(math.Pow(2, float64((card.Points - 1))))
		}
	}

	return points
}

type Card struct {
	WinningNumbers     []int
	ScratchcardNumbers []int
	Points             int
	Copies             int
}

func Part2(input string) int {

	lines := strings.Split(strings.TrimSpace(input), "\n")

	cards := createCards(lines)

	for idx, card := range cards {
		for i := 1; i <= card.Points && idx+i < len(cards); i++ {
			cards[idx+i].Copies += card.Copies
		}
	}

	totalScratchcards := 0
	for _, card := range cards {
		totalScratchcards += card.Copies
	}

	return totalScratchcards
}
