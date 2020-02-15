package main

import "fmt"

func fibonacci() func() int {
	a := 0
	b := 1
	return func() int {
		a, b = b, a+b
		return b
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 77; i++ {
		fmt.Println(f())
	}
}
