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

var cardValuesPart1 = [13]string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
var cardValuesPart2 = [13]string{"J", "2", "3", "4", "5", "6", "7", "8", "9", "T", "Q", "K", "A"}

type hand struct {
	cards []string
	bid   int
}

func getHands(lines []string) []hand {
	hands := []hand{}
	for _, line := range lines {
		split := strings.Split(line, " ")
		cards := strings.Split(split[0], "")
		bid, _ := strconv.Atoi(split[1])
		hands = append(hands, hand{cards: cards, bid: bid})
	}
	return hands
}

func cardCounts(cards []string) map[string]int {
	counts := make(map[string]int)
	for _, card := range cards {
		counts[card] += 1
	}
	return counts
}

func isXOfAKind(cards []string, x int) bool {
	counts := utils.Values(cardCounts(cards))
	return slices.Contains(counts, x)
}

func isFullHouse(cards []string) bool {
	counts := utils.Values(cardCounts(cards))
	if slices.Contains(counts, 3) && slices.Contains(counts, 2) {
		return true
	}
	return false
}

func isTwoPair(cards []string) bool {
	counts := utils.Values(cardCounts(cards))
	return len(utils.Filter(counts, func(value int) bool {
		return value == 2
	})) == 2
}

func getHandValue(cards []string) int {
	if isXOfAKind(cards, 5) {
		return 7
	}
	if isXOfAKind(cards, 4) {
		return 6
	}
	if isFullHouse(cards) {
		return 5
	}
	if isXOfAKind(cards, 3) {
		return 4
	}
	if isTwoPair(cards) {
		return 3
	}
	if isXOfAKind(cards, 2) {
		return 2
	}
	return 1
}

// Recursive function to replace a card in a hand with all possible replacements
func replaceAndCombine(cards []string, replace string, with []string) [][]string {
	idx := utils.IndexOf(len(cards), func(idx int) bool { return cards[idx] == replace })

	options := [][]string{}

	if idx == -1 {
		options = append(options, cards)
		return options
	}

	for _, card := range with {
		newCards := make([]string, len(cards))
		copy(newCards, cards)
		newCards[idx] = card
		nextOptions := replaceAndCombine(newCards, replace, with)
		options = append(options, nextOptions...)
	}

	return options
}

func getBestHandValue(cards []string) int {
	options := replaceAndCombine(cards, "J", []string{"A", "K", "Q", "T", "9", "8", "7", "6", "5", "4", "3", "2"})
	best := 0
	for _, option := range options {
		value := getHandValue(option)
		if value > best {
			best = value
		}
	}
	return best
}

func getCardValue(card string, part2 bool) int {
	cardValues := cardValuesPart1
	if part2 {
		cardValues = cardValuesPart2
	}
	return utils.IndexOf(len(cardValues), func(idx int) bool {
		return cardValues[idx] == card
	})
}

func Part1(input string) int {

	lines := strings.Split(strings.TrimSpace(input), "\n")

	hands := getHands(lines)

	sort.Slice(hands, func(i, j int) bool {
		a := getHandValue(hands[i].cards)
		b := getHandValue(hands[j].cards)

		if a == b {
			for x := 0; x < 5; x++ {
				cardValueA := getCardValue(hands[i].cards[x], false)
				cardValueB := getCardValue(hands[j].cards[x], false)
				if cardValueA < cardValueB {
					return true
				} else if cardValueA > cardValueB {
					return false
				}
			}
		}

		return a < b
	})

	totalWinnings := 0
	for idx, hand := range hands {
		totalWinnings = totalWinnings + (hand.bid * (idx + 1))
	}

	return totalWinnings
}

func Part2(input string) int {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	hands := getHands(lines)

	sort.Slice(hands, func(i, j int) bool {
		a := getBestHandValue(hands[i].cards)
		b := getBestHandValue(hands[j].cards)

		// If the hands are equal - check for the highest card
		if a == b {
			for x := 0; x < 5; x++ {
				cardValueA := getCardValue(hands[i].cards[x], true)
				cardValueB := getCardValue(hands[j].cards[x], true)
				if cardValueA < cardValueB {
					return true
				} else if cardValueA > cardValueB {
					return false
				}
			}
		}

		return a < b
	})

	totalWinnings := 0
	for idx, hand := range hands {
		totalWinnings = totalWinnings + (hand.bid * (idx + 1))
	}

	return totalWinnings
}
