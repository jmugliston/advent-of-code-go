package cli

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"golang.org/x/net/html"
)

const BASE_URL = "https://adventofcode.com"

var SESSION_COOKIE string
var USER_AGENT string

// InitialiseDay initialises the Advent of Code day for a given year and day.
// It creates the necessary folders, template files, and fetches the question and input for the specified day.
//
// Parameters:
//   - year: The year of the Advent of Code challenge.
//   - day: The day of the Advent of Code challenge.
//
// Example:
//
//	InitialiseDay("2021", "1")
func InitialiseDay(year string, day string) {
	dayPadded := day
	if len(day) == 1 {
		dayPadded = "0" + day
	}

	path := filepath.Join(".", year, "day"+dayPadded)

	logger.Info("Intialising day", "year", year, "day", day)

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		logger.Warn("Skipping template - folder already exists", "year", year, "day", day)
	} else {
		makeFolders(path)
		createTemplateFiles(path)
		FetchQuestion(year, day, path)
	}
	FetchInput(year, day, path)
}

// DownloadInput downloads the input for a given year and day.
// If day is set to 0, it downloads the input for all days of the year.
// The input is fetched and saved in the specified path.
//
// Parameters:
//   - year: The year for which to download the input.
//   - day: The day for which to download the input. If set to 0, it downloads the input for all days of the year.
//
// Example usage:
//
//	DownloadInput(2022, 1) // Downloads the input for year 2022, day 1
//	DownloadInput(2022, 0) // Downloads the input for all days of year 2022
func DownloadInput(year string, day string) {
	if day == "0" {
		for i := 1; i <= 25; i++ {
			nextDay := fmt.Sprintf("%d", i)
			nextDayPadded := fmt.Sprintf("%02d", i)
			path := filepath.Join(".", year, "day"+nextDayPadded)

			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			FetchInput(year, nextDay, path)
		}
	} else {
		dayPadded := day

		if len(day) == 1 {
			dayPadded = "0" + day
		}

		path := filepath.Join(".", year, "day"+dayPadded)

		FetchInput(year, day, path)
	}
}

// FetchQuestion fetches the question for a specific year and day from the Advent of Code website
// and saves it as a Markdown file in the specified path.
//
// Parameters:
//   - year: The year of the Advent of Code challenge.
//   - day: The day of the Advent of Code challenge.
//   - path: The path where the Markdown file will be saved.
//
// Example:
//
//	FetchQuestion("2021", "1", "/home/user/advent-of-code")
//
// Panics:
//   - If there is an error while fetching the question file.
//   - If the question for the specified day is not found (returns 404).
//   - If there is an error while parsing the HTML response.
//   - If the <article> element is not found in the HTML response.
//   - If there is an error while converting the HTML to Markdown.
func FetchQuestion(year string, day string, path string) {
	logger.Info("Downloading question for", "year", year, "day", day)

	url := fmt.Sprintf("%s/%s/day/%s", BASE_URL, year, day)

	resp, err := http.Get(url)

	if err != nil {
		logger.Error("Failed to fetch question file")
		panic(err)
	}

	if resp.StatusCode == 404 {
		logger.Error("404 - Could not find the question for this day")
		os.Exit(1)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		logger.Error(err)
	}

	articleElement := findArticleElement(doc)

	if articleElement == nil {
		logger.Error("Could not find the <article> element in the HTML")
		os.Exit(1)
	}

	articleHTML := renderNodeToHTML(articleElement)

	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(articleHTML)

	if err != nil {
		panic(err)
	}

	saveStringToFile(markdown, filepath.Join(path, "README.md"))
}

// FetchInput fetches the input file for a given year and day from the Advent of Code API and saves it to a specified path.
//
// Parameters:
//   - year: a string representing the year of the Advent of Code challenge
//   - day: a string representing the day of the Advent of Code challenge
//   - path: a string representing the path where the input file will be saved
//
// Example:
//
//	FetchInput("2021", "1", "/advent-of-code/2021/day01")
func FetchInput(year string, day string, path string) {

	if _, err := os.Stat(filepath.Join(path, "input", "input.txt")); err == nil {
		logger.Warn("Skipping download - input file already exists", "year", year, "day", day)
		return
	}

	logger.Info("Downloading input for", "year", year, "day", day)

	url := fmt.Sprintf("%s/%s/day/%s/input", BASE_URL, year, day)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", SESSION_COOKIE)
	req.Header.Set("User-Agent", USER_AGENT)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		logger.Errorf("Input file not found (%s)", url)
		os.Exit(1)
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	input := buf.String()

	saveStringToFile(input, filepath.Join(path, "input", "input.txt"))
}

// SolveDay solves the Advent of Code puzzle for a given year, day, and part.
// It executes the corresponding Go program and returns the output as a string.
// If the selected day does not exist, it logs an error and exits with code 1.
//
// Parameters:
//   - year: The year of the Advent of Code puzzle.
//   - day: The day of the Advent of Code puzzle.
//   - part: The part of the Advent of Code puzzle (1 or 2).
//
// Returns:
//   - string: The output of the Go program as a string.
func SolveDay(year string, day string, part string) string {
	logger.Info("Solving", "year", year, "day", day, "part", part)

	if len(day) == 1 {
		day = "0" + day
	}

	day = fmt.Sprintf("day%s", day)

	path := filepath.Join(".", year, day)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Error("Selected day does not exist")
		os.Exit(1)
	}

	out, err := exec.Command("go", "run", fmt.Sprintf("%s/main.go", path), "--part", part).Output()

	if err != nil {
		panic(err)
	}

	answer := strings.TrimSpace(string(out))

	fmt.Println(answer)

	return answer
}

// SubmitAnswer submits the answer for a given year, day, and part to the Advent of Code API.
// It posts the answer to the API and prints the response message.
//
// Parameters:
//   - year: The year of the Advent of Code challenge.
//   - day: The day of the Advent of Code challenge.
//   - part: The part of the Advent of Code challenge (1 or 2).
//
// Example:
//
//	SubmitAnswer("2021", "1", "1")
func SubmitAnswer(year string, day string, part string) {
	answer := strings.TrimSpace(SolveDay(year, day, part))

	url := fmt.Sprintf("%s/%s/day/%s/answer", BASE_URL, year, day)

	body := fmt.Sprintf("level=%s&answer=%s", part, answer)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Cookie", SESSION_COOKIE)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		logger.Error("Failed to submit answer")
		os.Exit(1)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)

	if err != nil {
		logger.Error("Could not parse response body")
		os.Exit(1)
	}

	articleElement := findArticleElement(doc)

	if articleElement == nil {
		logger.Error("Could not find the <article> element in the HTML")
		os.Exit(1)
	}

	text := extractNodeText(articleElement)

	if strings.Contains(text, "That's the right answer!") {
		fmt.Println("⭐ That's the right answer!")
	} else {
		lines := strings.Split(text, "  ")
		for _, line := range lines {
			if strings.Contains(line, "[") {
				line = strings.Split(line, " [")[0]
			}
			fmt.Println(line)
		}
	}
}
