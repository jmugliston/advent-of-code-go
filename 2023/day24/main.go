package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

	"github.com/atheius/aoc/parsing"
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
		fmt.Println(Part1(string(input), 200000000000000, 400000000000000))
	} else {
		fmt.Println(Part2(string(input)))
	}
}

var DIGIT_PATTERN = regexp.MustCompile(`-?\d+`)

type point struct {
	x int
	y int
	z int
}

type hailstone struct {
	position point
	velocity point
}

func parseHailstones(lines []string) []hailstone {
	hailstones := make([]hailstone, 0)

	for _, line := range lines {
		matches := DIGIT_PATTERN.FindAllString(line, -1)

		digits := make([]int, 0)
		for _, match := range matches {
			digit, err := strconv.Atoi(match)

			if err != nil {
				panic("Could not parse digit")
			}

			digits = append(digits, digit)
		}

		hailstone := hailstone{
			position: point{
				x: digits[0],
				y: digits[1],
				z: digits[2],
			},
			velocity: point{
				x: digits[3],
				y: digits[4],
				z: digits[5],
			},
		}

		hailstones = append(hailstones, hailstone)
	}

	return hailstones

}

func dot(a point, b point) int {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func cross(a point, b point) point {
	return point{
		x: a.y*b.z - a.z*b.y,
		y: a.z*b.x - a.x*b.z,
		z: a.x*b.y - a.y*b.x,
	}
}

func minus(a point, b point) point {
	return point{
		x: a.x - b.x,
		y: a.y - b.y,
		z: a.z - b.z,
	}
}

func plus(a point, b point) point {
	return point{
		x: a.x + b.x,
		y: a.y + b.y,
		z: a.z + b.z,
	}
}

func multiply(a point, b int) point {
	return point{
		x: a.x * b,
		y: a.y * b,
		z: a.z * b,
	}
}

func divide(a point, b int) point {
	return point{
		x: a.x / b,
		y: a.y / b,
		z: a.z / b,
	}
}

func findLinePlaneIntersection(p0 point, n point, stone hailstone) (point, int) {
	d := dot(minus(p0, stone.position), n) / dot(stone.velocity, n)
	return plus(stone.position, multiply(stone.velocity, d)), d
}

// Line intersect by Paul Bourke http://paulbourke.net/geometry/pointlineplane/
func lineIntersect(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, x4 int, y4 int) (bool, []float64) {

	// Check if none of the lines are of length 0
	if (x1 == x2 && y1 == y2) || (x3 == x4 && y3 == y4) {
		return false, nil
	}

	denominator := (y4-y3)*(x2-x1) - (x4-x3)*(y2-y1)

	// Lines are parallel
	if denominator == 0 {
		return false, nil
	}

	ua := float64(((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3))) / float64(denominator)

	// Return a object with the x and y coordinates of the intersection
	x := float64(x1) + ua*(float64(x2-x1))
	y := float64(y1) + ua*(float64(y2-y1))

	return true, []float64{x, y}
}

func xyWithinRange(x float64, y float64, min float64, max float64) bool {
	return x >= min && x <= max && y >= min && y <= max
}

func Part1(input string, minBound int, maxBound int) int {

	lines := parsing.ReadLines(input)

	hailstones := parseHailstones(lines)

	intersections := 0
	for i := 0; i < len(hailstones); i++ {
		for j := i + 1; j < len(hailstones); j++ {
			vx1 := hailstones[i].velocity.x
			vy1 := hailstones[i].velocity.y

			vx2 := hailstones[j].velocity.x
			vy2 := hailstones[j].velocity.y

			x1 := hailstones[i].position.x
			y1 := hailstones[i].position.y

			x2 := x1 + vx1
			y2 := y1 + vy1

			x3 := hailstones[j].position.x
			y3 := hailstones[j].position.y

			x4 := x3 + vx2
			y4 := y3 + vy2

			intersect, intersection := lineIntersect(x1, y1, x2, y2, x3, y3, x4, y4)

			if !intersect {
				continue
			}

			// Intersection must be in the future
			if (intersection[0]-float64(x1))/float64(vx1) < 0 || (intersection[0]-float64(x3))/float64(vx2) < 0 {
				continue
			}

			// Intersection must be within bounds
			if xyWithinRange(intersection[0], intersection[1], float64(minBound), float64(maxBound)) {
				intersections++
			}
		}
	}

	return intersections
}

// Credit to this comment on Reddit for the methodology of the solution
// https://www.reddit.com/r/adventofcode/comments/18q0kfc/comment/kes6ywf
func Part2(input string) int {

	lines := parsing.ReadLines(input)

	hailstones := parseHailstones(lines)

	// First hailstone used as a reference
	reference := hailstones[0]

	// Map hailstones relative to the first one
	relativeHailstones := make([]hailstone, 0)
	for _, stone := range hailstones {
		relativeHailstones = append(relativeHailstones, hailstone{
			position: minus(stone.position, reference.position),
			velocity: minus(stone.velocity, reference.velocity),
		})
	}

	// The first hailstone is at 0,0,0 and does not move
	// Which means the rock needs to pass through 0,0,0

	// The rock must intersect with the next hailstone somewhere in
	// the plane defined by the origin (0,0,0) and any two points on
	// the next hailstones trajectory

	hailstone1 := relativeHailstones[1]
	hailstone1position2 := plus(hailstone1.position, hailstone1.velocity)

	// Get the normal vector of hailstone 1
	n := cross(hailstone1.position, hailstone1position2)

	// Take two more hailstones and find the intersections their lines with the plane
	intersection1position, intersection1time := findLinePlaneIntersection(point{0, 0, 0}, n, relativeHailstones[2])
	intersection2position, intersection2time := findLinePlaneIntersection(point{0, 0, 0}, n, relativeHailstones[3])

	timeDiff := intersection2time - intersection1time

	// This is the relative rock velocity (velocity = distance / time)
	relativeRockVelocity := divide(minus(intersection2position, intersection1position), timeDiff)

	// This is the relative rock position (distance = velocity * time)
	relativeRockPosition := minus(intersection1position, multiply(relativeRockVelocity, intersection1time))

	// Convert back to absolute position
	rockPosition := plus(relativeRockPosition, reference.position)

	result := rockPosition.x + rockPosition.y + rockPosition.z

	return int(result)
}
