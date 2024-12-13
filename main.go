package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmugliston/aoc/cli"
	"github.com/joho/godotenv"
)

var Version string

var SESSION_TOKEN string

func main() {
	err := godotenv.Load(".env")

	SESSION_TOKEN = os.Getenv("SESSION_TOKEN")

	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	appVersion := Version
	if appVersion == "" {
		appVersion = "dev"
	}

	cli.VERSION = appVersion
	cli.SESSION_COOKIE = fmt.Sprintf("session=%s", SESSION_TOKEN)
	cli.USER_AGENT = fmt.Sprintf("github.com/jmugliston/aoc-go %s", appVersion)

	cli.Execute()
}
