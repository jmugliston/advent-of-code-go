package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/atheius/aoc/parsing"
	"github.com/juliangruber/go-intersect"
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
