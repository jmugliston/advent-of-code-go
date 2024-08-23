package xyz

type Coord struct {
	X int
	Y int
	Z int
}

func Dot(a Coord, b Coord) int {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross(a Coord, b Coord) Coord {
	return Coord{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func Minus(a Coord, b Coord) Coord {
	return Coord{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func Plus(a Coord, b Coord) Coord {
	return Coord{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

func Multiply(a Coord, b int) Coord {
	return Coord{
		X: a.X * b,
		Y: a.Y * b,
		Z: a.Z * b,
	}
}

func Divide(a Coord, b int) Coord {
	return Coord{
		X: a.X / b,
		Y: a.Y / b,
		Z: a.Z / b,
	}
}
