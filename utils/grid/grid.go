package grid

import (
	"fmt"
	"math"
	"strings"
)

type Grid [][]string

func (g Grid) Format(f fmt.State, c rune) {
	fmt.Fprintln(f, "")
	for _, line := range g {
		for _, char := range line {
			fmt.Fprintf(f, "%s ", char)
		}
		fmt.Fprintln(f)
	}
	fmt.Fprintln(f, "")
}

type Point struct {
	X int
	Y int
}

type Direction int

type PointWithDirection struct {
	X         int
	Y         int
	Direction Direction
}

const (
	North Direction = iota + 1
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

func (d Direction) String() string {
	return [...]string{"North", "NorthEast", "East", "SouthEast", "South", "SouthWest", "West", "NorthWest"}[d-1]
}

func (d Direction) EnumIndex() int {
	return int(d)
}

type DirectionalPoints struct {
	North     Point
	NorthEast Point
	East      Point
	SouthEast Point
	South     Point
	SouthWest Point
	West      Point
	NorthWest Point
}

func Parse(input string) Grid {
	var output Grid
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		output = append(output, strings.Split(line, ""))
	}
	return output
}

func SurroundingPoints(p Point) DirectionalPoints {
	return DirectionalPoints{
		North:     Point{X: p.X, Y: p.Y - 1},
		NorthEast: Point{X: p.X + 1, Y: p.Y - 1},
		East:      Point{X: p.X + 1, Y: p.Y},
		SouthEast: Point{X: p.X + 1, Y: p.Y + 1},
		South:     Point{X: p.X, Y: p.Y + 1},
		SouthWest: Point{X: p.X - 1, Y: p.Y + 1},
		West:      Point{X: p.X - 1, Y: p.Y},
		NorthWest: Point{X: p.X - 1, Y: p.Y - 1},
	}
}

func GetNextPointInDirection(p PointWithDirection) Point {
	switch p.Direction {
	case North:
		return Point{X: p.X, Y: p.Y - 1}
	case NorthEast:
		return Point{X: p.X + 1, Y: p.Y - 1}
	case East:
		return Point{X: p.X + 1, Y: p.Y}
	case SouthEast:
		return Point{X: p.X + 1, Y: p.Y + 1}
	case South:
		return Point{X: p.X, Y: p.Y + 1}
	case SouthWest:
		return Point{X: p.X - 1, Y: p.Y + 1}
	case West:
		return Point{X: p.X - 1, Y: p.Y}
	case NorthWest:
		return Point{X: p.X - 1, Y: p.Y - 1}
	default:
		return Point{}
	}
}

func IsPointInGrid(p Point, g Grid) bool {
	return p.Y >= 0 && p.Y < len(g) && p.X >= 0 && p.X < len(g[0])
}

func Transpose(a [][]string) [][]string {
	newArr := make([][]string, len(a[0]))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			newArr[j] = append(newArr[j], a[i][j])
		}
	}
	return newArr
}

func RotateClockwise(a [][]string) [][]string {
	newArr := make([][]string, len(a[0]))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			newArr[j] = append(newArr[j], a[len(a)-1-i][j])
		}
	}
	return newArr
}

func RotateCounterClockwise(a [][]string) [][]string {
	newArr := make([][]string, len(a[0]))
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			newArr[j] = append(newArr[j], a[i][len(a[0])-1-j])
		}
	}
	return newArr
}

func ManhattenDistance(p1 Point, p2 Point) int {
	return int(math.Abs(float64(p1.X)-float64(p2.X)) + math.Abs(float64(p1.Y)-float64((p2.Y))))
}

func Compare(g1 Grid, g2 Grid) bool {
	if len(g1) != len(g2) {
		return false
	}
	for i := range g1 {
		if len(g1[i]) != len(g2[i]) {
			return false
		}
		for j := range g1[i] {
			if g1[i][j] != g2[i][j] {
				return false
			}
		}
	}
	return true
}

func Copy(g Grid) Grid {
	newGrid := make(Grid, len(g))
	for i := range g {
		newGrid[i] = make([]string, len(g[i]))
		copy(newGrid[i], g[i])
	}
	return newGrid
}

func ToString(g Grid) string {
	var output strings.Builder
	for _, line := range g {
		output.WriteString(strings.Join(line, ""))
	}
	return output.String()
}
