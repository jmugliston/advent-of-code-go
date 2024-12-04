package grid

import (
	"fmt"
	"strconv"
	"strings"
)

func (g NumberGrid) Format(f fmt.State, c rune) {
	fmt.Fprintln(f, "")
	for _, line := range g {
		for _, num := range line {
			fmt.Fprintf(f, "%s ", fmt.Sprintf("%d", num))
		}
		fmt.Fprintln(f)
	}
	fmt.Fprintln(f, "")
}

func (g NumberGrid) ToString() string {
	var output strings.Builder
	for _, line := range g {
		var lineStr []string
		for _, num := range line {
			lineStr = append(lineStr, strconv.Itoa(num))
		}
		output.WriteString(strings.Join(lineStr, ""))
	}
	return output.String()
}

func (g NumberGrid) IsPointInGrid(p Point) bool {
	return p.Y >= 0 && p.Y < len(g) && p.X >= 0 && p.X < len(g[0])
}

func (g NumberGrid) GetPoint(p Point) int {
	return g[p.Y][p.X]
}

func (g NumberGrid) Transpose() NumberGrid {
	newArr := make(NumberGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[i][j])
		}
	}
	return newArr
}

func (g NumberGrid) RotateClockwise() NumberGrid {
	newArr := make(NumberGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[len(g)-1-i][j])
		}
	}
	return newArr
}

func (g NumberGrid) RotateCounterClockwise() NumberGrid {
	newArr := make(NumberGrid, len(g[0]))
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[0]); j++ {
			newArr[j] = append(newArr[j], g[i][len(g[0])-1-j])
		}
	}
	return newArr
}

func ConvertToStringGrid(grid NumberGrid) StringGrid {
	var newGrid StringGrid
	for _, line := range grid {
		var newLine []string
		for _, num := range line {
			newLine = append(newLine, strconv.Itoa(num))
		}
		newGrid = append(newGrid, newLine)
	}
	return newGrid
}

func (g NumberGrid) PrintPath(path []PointWithDirection) {
	grid := Copy(g)
	gridString := ConvertToStringGrid(grid)

	directionSymbols := map[Direction]string{
		North: "^",
		East:  ">",
		South: "v",
		West:  "<",
	}

	for _, step := range path {
		if symbol, exists := directionSymbols[step.Direction]; exists {
			gridString[step.Y][step.X] = symbol
		}
	}

	fmt.Println(gridString)
}
