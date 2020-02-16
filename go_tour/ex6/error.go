package main

import (
	"fmt"
)

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

	x1 := x / 2
	x2 := 0.0

	for {
		x2 = 0.5 * (x1 + x/x1)
		fmt.Println(x2)

		if (x1 - x2) < 0.00000000000001 {
			return x1, nil
		}
		x1 = x2
	}
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func main() {
	fmt.Println(Sqrt(10))
	fmt.Println(Sqrt(-2))
}
