package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	a := make([][]uint8, dy)
	for i := range a {
		a[i] = make([]uint8, dx)
	}
	for i := 0; i < dy; i++ {
		for j := 0; j < dx; j++ {
			a[i][j] = uint8((i + j) / 2)
		}
	}
	return a
}

func main() {
	pic.Show(Pic)
}
