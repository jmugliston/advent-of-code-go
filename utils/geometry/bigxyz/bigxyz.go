package bigxyz

import (
	"math/big"

	"github.com/atheius/aoc/bigInt"
)

type Coord struct {
	X *big.Int
	Y *big.Int
	Z *big.Int
}

func Dot(a Coord, b Coord) *big.Int {
	return bigInt.Add(bigInt.Add(bigInt.Mul(a.X, b.X), bigInt.Mul(a.Y, b.Y)), bigInt.Mul(a.Z, b.Z))
}

func Cross(a Coord, b Coord) Coord {
	return Coord{
		X: bigInt.Sub(bigInt.Mul(a.Y, b.Z), bigInt.Mul(a.Z, b.Y)),
		Y: bigInt.Sub(bigInt.Mul(a.Z, b.X), bigInt.Mul(a.X, b.Z)),
		Z: bigInt.Sub(bigInt.Mul(a.X, b.Y), bigInt.Mul(a.Y, b.X)),
	}
}

func Minus(a Coord, b Coord) Coord {
	return Coord{
		X: bigInt.Sub(a.X, b.X),
		Y: bigInt.Sub(a.Y, b.Y),
		Z: bigInt.Sub(a.Z, b.Z),
	}
}

func Plus(a Coord, b Coord) Coord {
	return Coord{
		X: bigInt.Add(a.X, b.X),
		Y: bigInt.Add(a.Y, b.Y),
		Z: bigInt.Add(a.Z, b.Z),
	}
}

func Multiply(a Coord, b *big.Int) Coord {
	return Coord{
		X: bigInt.Mul(a.X, b),
		Y: bigInt.Mul(a.Y, b),
		Z: bigInt.Mul(a.Z, b),
	}
}

func Divide(a Coord, b *big.Int) Coord {
	return Coord{
		X: bigInt.Div(a.X, b),
		Y: bigInt.Div(a.Y, b),
		Z: bigInt.Div(a.Z, b),
	}
}
