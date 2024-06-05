package main

import (
	"fmt"
	"os"

	"github.com/atheius/aoc/grid"
	"github.com/atheius/aoc/utils"
	"gonum.org/v1/gonum/stat/combin"
)

func main() {
	input, err := os.ReadFile("./input/input.txt")

	if err != nil {
		panic("Couldn't find the input file!")
	}

	inputString := string(input)

	part1Answer := Part1(inputString)
	fmt.Println(part1Answer)

	part2Answer := Part2(inputString, 1000000)
	fmt.Println(part2Answer)
}

func getEmptyRows(image [][]string) []int {

	var emptyRows []int

	for y, row := range image {
		empty := true
		for _, cell := range row {
			if cell == "#" {
				empty = false
				break
			}
		}
		if empty {
			emptyRows = append(emptyRows, y)
		}
	}

	return emptyRows
}

func getGalaxies(image [][]string) []grid.Point {

	var galaxies []grid.Point

	for y, row := range image {
		for x, cell := range row {
			if cell == "#" {
				galaxies = append(galaxies, grid.Point{X: x, Y: y})
			}
		}
	}

	return galaxies
}

func getScaledGalaxy(galaxy grid.Point, emptyRows []int, emptyCols []int, scale int) grid.Point {
	numEmptyRows := len(utils.Filter(emptyRows, func(row int) bool {
		return row < galaxy.Y
	}))

	numEmptyCols := len(utils.Filter(emptyCols, func(row int) bool {
		return row < galaxy.X
	}))

	scaledPoint := grid.Point{
		X: galaxy.X + numEmptyCols*(scale-1),
		Y: galaxy.Y + numEmptyRows*(scale-1),
	}

	return scaledPoint
}

func getScaledGalaxies(galaxies []grid.Point, emptyRows []int, emptyCols []int, scale int) []grid.Point {
	var scaledGalaxies []grid.Point
	for _, galaxy := range galaxies {
		scaledGalaxies = append(scaledGalaxies, getScaledGalaxy(galaxy, emptyRows, emptyCols, scale))
	}
	return scaledGalaxies
}

func getDistances(galaxies []grid.Point) []int {
	combinations := combin.Combinations(len(galaxies), 2)

	var distances []int
	for _, combination := range combinations {
		distance := grid.ManhattenDistance(galaxies[combination[0]], galaxies[combination[1]])
		distances = append(distances, distance)
	}

	return distances
}

func Part1(input string) int {

	image := grid.Parse(input)

	galaxies := getGalaxies(image)

	emptyRows := getEmptyRows(image)
	emptyCols := getEmptyRows(image.Transpose())

	scaledGalaxies := getScaledGalaxies(galaxies, emptyRows, emptyCols, 2)

	distances := getDistances(scaledGalaxies)

	sum := 0
	for _, distance := range distances {
		sum += distance
	}

	return sum
}

func Part2(input string, scale int) int {
	image := grid.Parse(input)

	galaxies := getGalaxies(image)

	emptyRows := getEmptyRows(image)
	emptyCols := getEmptyRows(image.Transpose())

	scaledGalaxies := getScaledGalaxies(galaxies, emptyRows, emptyCols, scale)

	distances := getDistances(scaledGalaxies)

	sum := 0
	for _, distance := range distances {
		sum += distance
	}

	return sum
}
