package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/jmugliston/aoc/grid"
	"github.com/jmugliston/aoc/parsing"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "input", "input.txt")

	input, err := os.ReadFile(path)

	if err != nil {
		panic("Could not find the input file")
	}

	if *partFlag == "1" {
		fmt.Println(Part1(string(input), false))
	} else {
		fmt.Println(Part2(string(input), true))
	}
}

type Robot struct {
	point grid.Point
	vx    int
	vy    int
}

func parseRobots(input string) []*Robot {
	robots := make([]*Robot, 0)

	lines := parsing.ReadLines(input)

	for _, line := range lines {
		split := strings.Split(line, " ")

		position := strings.Split(strings.Split(split[0], "=")[1], ",")
		velocity := strings.Split(strings.Split(split[1], "=")[1], ",")

		x, _ := strconv.Atoi(position[0])
		y, _ := strconv.Atoi(position[1])
		vx, _ := strconv.Atoi(velocity[0])
		vy, _ := strconv.Atoi(velocity[1])

		robots = append(robots, &Robot{grid.Point{X: x, Y: y}, vx, vy})
	}

	return robots
}

func (robot *Robot) move(boundary []int) Robot {
	boundaryX := boundary[0]
	boundaryY := boundary[1]

	newPosX := (robot.point.X + robot.vx) % boundaryX
	newPosY := (robot.point.Y + robot.vy) % boundaryY

	if newPosX < 0 {
		newPosX = boundaryX + newPosX
	}

	if newPosX > boundaryX {
		newPosX = newPosX - boundaryX
	}

	if newPosY < 0 {
		newPosY = boundaryY + newPosY
	}

	if newPosY > boundaryY {
		newPosY = newPosY - boundaryY
	}

	robot.point.X = newPosX
	robot.point.Y = newPosY

	return *robot
}

func (robot *Robot) getQuadrant(boundary []int) string {
	boundaryX := boundary[0]
	boundaryY := boundary[1]

	if robot.point.Y == boundaryY/2 || robot.point.X == boundaryX/2 {
		return "None"
	}

	if robot.point.X >= 0 && robot.point.X < boundaryX/2 {
		if robot.point.Y >= 0 && robot.point.Y < boundaryY/2 {
			return "NW"
		} else {
			return "SW"
		}
	} else {
		if robot.point.Y >= 0 && robot.point.Y < boundaryY/2 {
			return "NE"
		} else {
			return "SE"
		}
	}
}

func writeMapToFile(robotMap grid.StringGrid) {
	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "solution.txt")

	file, err := os.Create(path)

	if err != nil {
		panic("Could not create file")
	}

	defer file.Close()

	for _, row := range robotMap {
		for _, point := range row {
			file.WriteString(point)
		}
		file.WriteString("\n")
	}
}

func Part1(input string, example bool) int {
	robots := parseRobots(input)

	height := 103
	width := 101

	if example {
		height = 7
		width = 11
	}

	boundary := []int{width, height}

	for i := 0; i < 100; i++ {
		for _, robot := range robots {
			robot.move(boundary)
		}
	}

	quadrantMap := map[string]int{}
	for _, robot := range robots {
		quadrantMap[robot.getQuadrant(boundary)] += 1
	}

	return quadrantMap["NW"] * quadrantMap["NE"] * quadrantMap["SW"] * quadrantMap["SE"]
}

func Part2(input string, writeFile bool) int {
	robots := parseRobots(input)

	height := 103
	width := 101

	robotMap := grid.StringGrid{}

	// Intialise a robot map
	for i := 0; i < height; i++ {
		robotMap = append(robotMap, make([]string, width))
		for j := 0; j < width; j++ {
			robotMap[i][j] = "."
		}
	}

	boundary := []int{width, height}

	for i := 0; i < 10000; i++ {
		for _, robot := range robots {
			robot.move(boundary)
			robotMap.SetPoint(robot.point, "#")
		}

		for _, row := range robotMap {
			count := 0
			for _, point := range row {
				if point != "#" {
					count = 0
					continue
				}

				count++

				if count > 20 {
					// See solution.txt for the easter egg
					if writeFile {
						writeMapToFile(robotMap)
					}
					return i + 1
				}
			}
		}

		// Reset the robot map
		for _, robot := range robots {
			robotMap.SetPoint(robot.point, ".")
		}
	}

	return -1
}
