package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	x1 := x / 2
	x2 := 0.0

	for {
		x2 = 0.5 * (x1 + x/x1)
		fmt.Println(x2)

		if (x1 - x2) < 0.00000000000001 {
			return x1
		}
		x1 = x2
	}
}

func main() {
	fmt.Println(Sqrt(10))
}
