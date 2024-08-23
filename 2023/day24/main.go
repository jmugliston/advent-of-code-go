package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

	"github.com/atheius/aoc/bigInt"
	"github.com/atheius/aoc/bigxyz"
	"github.com/atheius/aoc/parsing"
	"github.com/atheius/aoc/xyz"
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

type hailstone struct {
	position xyz.Coord
	velocity xyz.Coord
}

type hailstoneBig struct {
	position bigxyz.Coord
	velocity bigxyz.Coord
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
			position: xyz.Coord{
				X: digits[0],
				Y: digits[1],
				Z: digits[2],
			},
			velocity: xyz.Coord{
				X: digits[3],
				Y: digits[4],
				Z: digits[5],
			},
		}

		hailstones = append(hailstones, hailstone)
	}

	return hailstones

}

func findLinePlaneIntersection(p0 bigxyz.Coord, n bigxyz.Coord, stone hailstoneBig) (bigxyz.Coord, *big.Int) {
	d := bigInt.Div(
		bigxyz.Dot(bigxyz.Minus(p0, stone.position), n),
		bigxyz.Dot(stone.velocity, n),
	)
	return bigxyz.Plus(stone.position, bigxyz.Multiply(stone.velocity, d)), d
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
			vx1 := hailstones[i].velocity.X
			vy1 := hailstones[i].velocity.Y

			vx2 := hailstones[j].velocity.X
			vy2 := hailstones[j].velocity.Y

			x1 := hailstones[i].position.X
			y1 := hailstones[i].position.Y

			x2 := x1 + vx1
			y2 := y1 + vy1

			x3 := hailstones[j].position.X
			y3 := hailstones[j].position.Y

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

func convertToCoordBig(coord xyz.Coord) bigxyz.Coord {
	return bigxyz.Coord{
		X: big.NewInt(int64(coord.X)),
		Y: big.NewInt(int64(coord.Y)),
		Z: big.NewInt(int64(coord.Z)),
	}
}

// Credit to this comment on Reddit for the methodology of the solution
// https://www.reddit.com/r/adventofcode/comments/18q0kfc/comment/kes6ywf
func Part2(input string) int {

	lines := parsing.ReadLines(input)

	hailstones := parseHailstones(lines)

	// First hailstone used as a reference
	reference := hailstones[0]

	// Map hailstones relative to the first one
	relativeHailstones := make([]hailstoneBig, 0)
	for _, stone := range hailstones {
		relativeHailstones = append(relativeHailstones, hailstoneBig{
			// We have to use big ints for this solution because the numbers are too large!
			position: convertToCoordBig(xyz.Minus(stone.position, reference.position)),
			velocity: convertToCoordBig(xyz.Minus(stone.velocity, reference.velocity)),
		})
	}

	// The first hailstone is at 0,0,0 and does not move
	// Which means the rock needs to pass through 0,0,0

	// The rock must intersect with the next hailstone somewhere in
	// the plane defined by the origin (0,0,0) and any two points on
	// the next hailstones trajectory

	// Get the normal vector of hailstone 1
	hailstone1 := relativeHailstones[1]
	hailstone1position2 := bigxyz.Plus(hailstone1.position, hailstone1.velocity)
	n := bigxyz.Cross(hailstone1.position, hailstone1position2)

	// Take two more hailstones and find the intersections their lines with the plane
	intersection1position, intersection1time := findLinePlaneIntersection(
		bigxyz.Coord{X: big.NewInt(0), Y: big.NewInt(0), Z: big.NewInt(0)},
		n,
		relativeHailstones[2],
	)
	intersection2position, intersection2time := findLinePlaneIntersection(
		bigxyz.Coord{X: big.NewInt(0), Y: big.NewInt(0), Z: big.NewInt(0)},
		n,
		relativeHailstones[3],
	)

	timeDiff := bigInt.Sub(intersection2time, intersection1time)

	// This is the relative rock velocity (velocity = distance / time)
	relativeRockVelocity := bigxyz.Divide(bigxyz.Minus(intersection2position, intersection1position), timeDiff)

	// This is the relative rock position (distance = velocity * time)
	relativeRockPosition := bigxyz.Minus(intersection1position, bigxyz.Multiply(relativeRockVelocity, intersection1time))

	// Convert back to absolute position
	rockPosition := bigxyz.Plus(relativeRockPosition, convertToCoordBig(reference.position))

	result := bigInt.Add(bigInt.Add(rockPosition.X, rockPosition.Y), rockPosition.Z)

	return int(result.Int64())
}
