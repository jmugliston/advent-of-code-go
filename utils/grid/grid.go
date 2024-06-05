package grid

import (
	"math"
	"strconv"
	"strings"
)

type Grid[N string | int] interface {
	~[][]N
}

type StringGrid [][]string
type NumberGrid [][]int

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

func Parse(input string) StringGrid {
	var output StringGrid
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		split := strings.Split(line, "")
		output = append(output, split)
	}
	return output
}

func ParseNumbers(input string) NumberGrid {
	var output NumberGrid
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		var split []int
		for _, char := range strings.Split(line, "") {
			num, err := strconv.Atoi(char)
			if err != nil {
				panic("Can't parse grid character as number")
			}
			split = append(split, num)
		}
		output = append(output, split)
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

func ManhattenDistance(p1 Point, p2 Point) int {
	return int(math.Abs(float64(p1.X)-float64(p2.X)) + math.Abs(float64(p1.Y)-float64((p2.Y))))
}

func Compare[G Grid[N], N string | int](g1 G, g2 G) bool {
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

func Copy[G Grid[N], N string | int](g G) G {
	newGrid := make(G, len(g))
	for i := range g {
		newGrid[i] = make([]N, len(g[i]))
		copy(newGrid[i], g[i])
	}
	return newGrid
}
