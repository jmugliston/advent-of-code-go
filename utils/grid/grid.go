package grid

import (
	"math"
	"strings"
)

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

func Transpose(a [][]string) [][]string {
	newArr := make([][]string, len(a))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			newArr[j] = append(newArr[j], a[i][j])
		}
	}
	return newArr
}

func ManhattenDistance(p1 Point, p2 Point) int {
	return int(math.Abs(float64(p1.X)-float64(p2.X)) + math.Abs(float64(p1.Y)-float64((p2.Y))))
}
