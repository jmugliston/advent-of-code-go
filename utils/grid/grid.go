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

func ReadGrid(input string) Grid {
	var output Grid
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
