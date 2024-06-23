package cli

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

// An interactive command-line interface for Advent of Code.
// It presents a menu of options to the user, such as initialising a new day, downloading puzzle input,
// solving a puzzle, or submitting an answer. The user can select the desired option by entering the
// corresponding number or by using arrow keys to navigate the menu.
func Interactive() {
	fmt.Println("\n🎄🎄🎄 Advent of Code 🎄🎄🎄")
	fmt.Println("----------------------------")

	currentYear, currentMonth, _ := time.Now().Date()

	defaultYear := currentYear
	if currentMonth != 12 {
		defaultYear = defaultYear - 1
	}

	years := []string{}
	for i := defaultYear; i >= 2015; i-- {
		years = append(years, fmt.Sprintf("%d", i))
	}

	optionPrompt := promptui.Select{
		Label: "What would you like to do?",
		Items: []string{"Initialise", "Download", "Solve", "Submit", "Exit"},
	}

	_, option, err := optionPrompt.Run()

	if err != nil {
		return
	}

	if option == "Exit" {
		return
	}

	yearPrompt := promptui.Select{
		Label:     "Which year?",
		Items:     years,
		CursorPos: 0,
	}

	_, year, err := yearPrompt.Run()

	if err != nil {
		return
	}

	dayPrompt := promptui.Prompt{
		Label:    "Day",
		Validate: validateDay,
	}

	day, err := dayPrompt.Run()

	if err != nil {
		return
	}

	if option == "Initialise" {
		InitialiseDay(year, day)
		return
	}

	if option == "Download" {
		DownloadInput(year, day)
		return
	}

	partPrompt := promptui.Select{
		Label: "Which part?",
		Items: []string{"1", "2"},
	}

	_, part, err := partPrompt.Run()

	if err != nil {
		return
	}

	if option == "Solve" {
		SolveDay(year, day, part)
		return
	}

	if option == "Submit" {
		SubmitAnswer(year, day, part)
		return
	}

}
