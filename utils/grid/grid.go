package grid

import "strings"

type Point struct {
	X int
	Y int
}

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d-1]
}

func (d Direction) EnumIndex() int {
	return int(d)
}

type DirectionalPoints struct {
	North Point
	East  Point
	South Point
	West  Point
}

func ReadGrid(input string) [][]string {
	var output [][]string
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		output = append(output, strings.Split(line, ""))
	}
	return output
}

func SurroundingPoints(p Point) DirectionalPoints {
	return DirectionalPoints{
		North: Point{X: p.X, Y: p.Y - 1},
		East:  Point{X: p.X + 1, Y: p.Y},
		South: Point{X: p.X, Y: p.Y + 1},
		West:  Point{X: p.X - 1, Y: p.Y},
	}
}
