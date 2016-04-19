package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image/color"
)

func main() {
	red := color.NRGBA{255, 0, 0, 255}
	dst := imaging.New(100, 100, color.NRGBA{0, 0, 0, 0})
	dst.SetNRGBA(52, 41, red)
	dst = imaging.FlipV(dst)
	err := imaging.Save(dst, "test.jpeg")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
