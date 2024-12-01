package utils

import (
	"math"
)

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func IndexOf(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func Sum(values []int) int {
	sum := 0
	for _, value := range values {
		sum += value
	}
	return sum
}

func Product(values []int) int {
	product := 1
	for _, value := range values {
		product *= value
	}
	return product
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(values []int) int {
	lcm := values[0]
	for _, value := range values[1:] {
		lcm = int(math.Abs(float64(lcm*value)) / float64(GCD(lcm, value)))
	}
	return lcm
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
