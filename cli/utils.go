package cli

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jmugliston/aoc/templates"
	"golang.org/x/net/html"
)

func getPaddedDay(day string) string {
	dayPadded := day
	if len(day) == 1 {
		dayPadded = "0" + day
	}
	return dayPadded
}

func validateDay(input string) error {
	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > 25 {
		return errors.New("invalid day (must be between 1 and 25)")
	}
	return nil
}

func createTemplateFiles(path string) {
	saveStringToFile(templates.MainTemplate, filepath.Join(path, "main.go"))
	saveStringToFile(templates.TestTemplate, filepath.Join(path, "main_test.go"))
	saveStringToFile("", filepath.Join(path, "input", "example.txt"))
}

func makeFolders(path string) {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(path, "input"), os.ModePerm)

	if err != nil {
		panic(err)
	}
}

func saveStringToFile(data string, path string) {
	file, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = file.WriteString(data)

	if err != nil {
		log.Fatal(err)
	}
}

func findArticleElements(n *html.Node) []*html.Node {
	var articles []*html.Node
	if n.Type == html.ElementNode && n.Data == "article" {
		articles = append(articles, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		articles = append(articles, findArticleElements(c)...)
	}
	return articles
}

func getQuestionHTML(n *html.Node) string {
	root := &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}

	articles := findArticleElements(n)

	for _, article := range articles {
		if article.Parent != nil {
			article.Parent.RemoveChild(article)
		}
		root.AppendChild(article)
	}

	var buf bytes.Buffer
	html.Render(&buf, root)
	return buf.String()
}

func extractNodeText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return buf.String()
}
