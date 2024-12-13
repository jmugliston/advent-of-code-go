package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/parsing"
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

type PartName string

const (
	X PartName = "x"
	M PartName = "m"
	A PartName = "a"
	s PartName = "s"
)

type Rule struct {
	Input    PartName
	Operator string
	Value    int
	Output   string
}

type Workflow struct {
	Name   string
	Rules  []Rule
	Output string
}

type PartRating struct {
	X int
	M int
	A int
	S int
}

// Types for part 2

type Range struct {
	Min int
	Max int
}

type PartRange struct {
	X Range
	M Range
	A Range
	S Range
}

type QueueItem struct {
	State string
	Workflow
	PartRange
}

var operatorRegex = regexp.MustCompile("<|>")

func ParseRule(rule string) Rule {

	ruleSplit := strings.Split(rule, ":")

	output := ruleSplit[1]

	operator := "<"
	if strings.Contains(ruleSplit[0], ">") {
		operator = ">"
	}

	actionSplit := operatorRegex.Split(ruleSplit[0], -1)

	input := actionSplit[0]
	value, _ := strconv.Atoi(actionSplit[1])

	return Rule{
		Input:    PartName(input),
		Operator: operator,
		Value:    value,
		Output:   output,
	}
}

func ParseWorkflow(rawWorkflow string) Workflow {

	split := strings.Split(rawWorkflow, "{")
	name := split[0]
	rules := strings.Trim(split[1], "}")

	splitRules := strings.Split(rules, ",")
	output := splitRules[len(splitRules)-1]

	parsedRules := []Rule{}
	for _, rule := range splitRules[:len(splitRules)-1] {
		nextRule := ParseRule(rule)
		parsedRules = append(parsedRules, nextRule)
	}

	return Workflow{
		Name:   name,
		Rules:  parsedRules,
		Output: output,
	}

}

func ParsePartRating(partRating string) PartRating {
	partRatingSplit := strings.Split(strings.Trim(partRating, "{}"), ",")

	x, _ := strconv.Atoi(partRatingSplit[0][2:])
	m, _ := strconv.Atoi(partRatingSplit[1][2:])
	a, _ := strconv.Atoi(partRatingSplit[2][2:])
	s, _ := strconv.Atoi(partRatingSplit[3][2:])

	return PartRating{
		X: x,
		M: m,
		A: a,
		S: s,
	}

}

func evaluateOperation(partValue int, operator string, ruleValue int) bool {
	switch operator {
	case "<":
		return partValue < ruleValue
	case ">":
		return partValue > ruleValue
	}
	panic("Invalid operator")
}

func evaluateRule(partRating PartRating, rule Rule) bool {
	switch rule.Input {
	case "x":
		return evaluateOperation(partRating.X, rule.Operator, rule.Value)
	case "m":
		return evaluateOperation(partRating.M, rule.Operator, rule.Value)
	case "a":
		return evaluateOperation(partRating.A, rule.Operator, rule.Value)
	case "s":
		return evaluateOperation(partRating.S, rule.Operator, rule.Value)
	}
	panic("Invalid rule input")
}

func Part1(input string) int {

	blocks := strings.Split(input, "\n\n")

	workflowMap := make(map[string]Workflow)
	for _, rawWorkflow := range parsing.ReadLines(blocks[0]) {
		parsedWorkflow := ParseWorkflow(rawWorkflow)
		workflowMap[parsedWorkflow.Name] = parsedWorkflow
	}

	var partRatings []PartRating
	for _, partRating := range parsing.ReadLines(blocks[1]) {
		partRatings = append(partRatings, ParsePartRating(partRating))
	}

	var acceptedParts []PartRating

	for _, partRating := range partRatings {
		currentPoint := "in"
		for currentPoint != "A" && currentPoint != "R" {
			workflow := workflowMap[currentPoint]
			next := workflow.Output
			for _, ruleToCheck := range workflow.Rules {
				check := evaluateRule(partRating, ruleToCheck)
				if check {
					next = ruleToCheck.Output
					break
				}
			}
			currentPoint = next
		}
		if currentPoint == "A" {
			acceptedParts = append(acceptedParts, PartRating{X: partRating.X, M: partRating.M, A: partRating.A, S: partRating.S})
		}
	}

	total := 0
	for _, partRating := range acceptedParts {
		total += partRating.X + partRating.M + partRating.A + partRating.S
	}

	return total
}

func Part2(input string) int {

	blocks := strings.Split(input, "\n\n")

	workflowMap := make(map[string]Workflow)
	for _, rawWorkflow := range parsing.ReadLines(blocks[0]) {
		parsedWorkflow := ParseWorkflow(rawWorkflow)
		workflowMap[parsedWorkflow.Name] = parsedWorkflow
	}

	queue := []QueueItem{
		{
			State:    "in",
			Workflow: workflowMap["in"],
			PartRange: PartRange{
				X: Range{Min: 1, Max: 4000},
				M: Range{Min: 1, Max: 4000},
				A: Range{Min: 1, Max: 4000},
				S: Range{Min: 1, Max: 4000},
			},
		},
	}

	acceptedRanges := []PartRange{}

	for len(queue) > 0 {

		currentItem := queue[0]
		queue = queue[1:]

		if currentItem.State == "A" {
			acceptedRanges = append(acceptedRanges, currentItem.PartRange)
			continue
		}

		if currentItem.State == "R" {
			continue
		}

		nextDefaultPartRange := PartRange{
			X: Range{Min: currentItem.PartRange.X.Min, Max: currentItem.PartRange.X.Max},
			M: Range{Min: currentItem.PartRange.M.Min, Max: currentItem.PartRange.M.Max},
			A: Range{Min: currentItem.PartRange.A.Min, Max: currentItem.PartRange.A.Max},
			S: Range{Min: currentItem.PartRange.S.Min, Max: currentItem.PartRange.S.Max},
		}

		for _, ruleToCheck := range currentItem.Workflow.Rules {

			if ruleToCheck.Operator == "<" {
				nextPartRange := nextDefaultPartRange

				value := ruleToCheck.Value
				if ruleToCheck.Input == "x" {
					nextPartRange.X.Max = value - 1
					nextDefaultPartRange.X.Min = value
				}

				if ruleToCheck.Input == "m" {
					nextPartRange.M.Max = value - 1
					nextDefaultPartRange.M.Min = value
				}

				if ruleToCheck.Input == "a" {
					nextPartRange.A.Max = value - 1
					nextDefaultPartRange.A.Min = value
				}

				if ruleToCheck.Input == "s" {
					nextPartRange.S.Max = value - 1
					nextDefaultPartRange.S.Min = value
				}

				nextState := ruleToCheck.Output
				nextItem := QueueItem{
					State:     nextState,
					Workflow:  workflowMap[nextState],
					PartRange: nextPartRange,
				}

				queue = append(queue, nextItem)
			}

			if ruleToCheck.Operator == ">" {
				nextPartRange := nextDefaultPartRange

				value := ruleToCheck.Value
				if ruleToCheck.Input == "x" {
					nextPartRange.X.Min = value + 1
					nextDefaultPartRange.X.Max = value
				}

				if ruleToCheck.Input == "m" {
					nextPartRange.M.Min = value + 1
					nextDefaultPartRange.M.Max = value
				}

				if ruleToCheck.Input == "a" {
					nextPartRange.A.Min = value + 1
					nextDefaultPartRange.A.Max = value
				}

				if ruleToCheck.Input == "s" {
					nextPartRange.S.Min = value + 1
					nextDefaultPartRange.S.Max = value
				}

				nextState := ruleToCheck.Output
				nextItem := QueueItem{
					State:     nextState,
					Workflow:  workflowMap[nextState],
					PartRange: nextPartRange,
				}

				queue = append(queue, nextItem)
			}
		}

		nextState := currentItem.Workflow.Output
		queue = append(queue, QueueItem{
			State:     nextState,
			Workflow:  workflowMap[nextState],
			PartRange: nextDefaultPartRange,
		})
	}

	total := 0
	for _, acceptedRange := range acceptedRanges {
		total += (acceptedRange.X.Max - acceptedRange.X.Min + 1) *
			(acceptedRange.M.Max - acceptedRange.M.Min + 1) *
			(acceptedRange.A.Max - acceptedRange.A.Min + 1) *
			(acceptedRange.S.Max - acceptedRange.S.Min + 1)
	}

	return total
}
