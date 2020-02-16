package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

type Image struct {
	Width  int
	Height int
	Color  int
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.Width, i.Height)
}

func (i Image) At(x, y int) color.Color {
	return color.RGBA{uint8(i.Color + x), uint8(i.Color + y), 255, 255}
}

func main() {
	m := Image{100, 100, 255}
	pic.ShowImage(m)
}
