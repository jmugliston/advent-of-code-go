module github.com/atheius/aoc

go 1.22.2

replace github.com/atheius/aoc/utils => ./utils/general

replace github.com/atheius/aoc/grid => ./utils/grid

replace github.com/atheius/aoc/parsing => ./utils/parsing

require (
	github.com/atheius/aoc/grid v0.0.0-00010101000000-000000000000
	github.com/atheius/aoc/parsing v0.0.0-00010101000000-000000000000
	github.com/atheius/aoc/utils v0.0.0-00010101000000-000000000000
	github.com/juliangruber/go-intersect v1.1.0
)

require gonum.org/v1/gonum v0.15.0
