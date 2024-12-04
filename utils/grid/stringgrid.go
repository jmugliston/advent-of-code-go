package grid

import (
	"fmt"
	"strings"
)

func (g StringGrid) Format(f fmt.State, c rune) {
	fmt.Fprintln(f, "")
	for _, line := range g {
		for _, char := range line {
			fmt.Fprintf(f, "%s ", char)
		}
		fmt.Fprintln(f)
	}
	fmt.Fprintln(f, "")
}

func (g StringGrid) ToString() string {
	var output strings.Builder
	for _, line := range g {
		output.WriteString(strings.Join(line, ""))
	}
	return output.String()
}

func (g StringGrid) IsPointInGrid(p Point) bool {
	return p.Y >= 0 && p.Y < len(g) && p.X >= 0 && p.X < len(g[0])
}

func (g StringGrid) GetPoint(p Point) string {
	if !g.IsPointInGrid(p) {
		return ""
	}
	return g[p.Y][p.X]
}

func (g StringGrid) GetPoints(p []Point) []string {
	points := []string{}
	for _, point := range p {
		points = append(points, g.GetPoint(point))
	}
	return points
}

func (g StringGrid) Find(char string) Point {
	for y, line := range g {
		for x, c := range line {
			if c == char {
				return Point{X: x, Y: y}
			}
		}
	}
	return Point{}
}

func (g StringGrid) FindAll(char string) []Point {
	points := []Point{}
	for y, line := range g {
		for x, c := range line {
			if c == char {
				points = append(points, Point{X: x, Y: y})
			}
		}
	}
	return points
}

func (g StringGrid) Transpose() StringGrid {
	newArr := make(StringGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[i][j])
		}
	}
	return newArr
}

func (g StringGrid) RotateClockwise() StringGrid {
	newArr := make(StringGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[len(g)-1-i][j])
		}
	}
	return newArr
}

func (g StringGrid) RotateCounterClockwise() StringGrid {
	newArr := make(StringGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[i][len(g[0])-1-j])
		}
	}
	return newArr
}

func (g StringGrid) PrintPath(path []PointWithDirection) {
	grid := Copy(g)

	directionSymbols := map[Direction]string{
		North: "^",
		East:  ">",
		South: "v",
		West:  "<",
	}

	for _, step := range path {
		if symbol, exists := directionSymbols[step.Direction]; exists {
			grid[step.Y][step.X] = symbol
		}
	}

	fmt.Println(grid)
}
