module github.com/atheius/aoc

go 1.22.2

replace github.com/atheius/aoc/cli => ./cli

replace github.com/atheius/aoc/utils => ./utils/general

replace github.com/atheius/aoc/grid => ./utils/grid

replace github.com/atheius/aoc/parsing => ./utils/parsing

require (
	github.com/atheius/aoc/grid v0.0.0-00010101000000-000000000000
	github.com/atheius/aoc/parsing v0.0.0-00010101000000-000000000000
	github.com/atheius/aoc/utils v0.0.0-00010101000000-000000000000
	github.com/juliangruber/go-intersect v1.1.0
)

require (
	github.com/JohannesKaufmann/html-to-markdown v1.6.0
	github.com/charmbracelet/log v0.4.0
	github.com/joho/godotenv v1.5.1
	github.com/manifoldco/promptui v0.9.0
	github.com/spf13/cobra v1.8.1
	golang.org/x/net v0.26.0
	gonum.org/v1/gonum v0.15.0
)

require (
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/lipgloss v0.10.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/sys v0.21.0 // indirect
)
