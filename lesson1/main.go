package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func line(x0, y0, x1, y1 int, img *image.NRGBA, color color.NRGBA) {

	if (x0 == x1) && (y0 == y1) {
		// point
		img.SetNRGBA(x0, y0, color)
	} else if x0 == x1 {
		// vertical line
		ymin := min(y0, y1)
		ymax := max(y0, y1)
		for y := ymin; y <= ymax; y++ {
			img.SetNRGBA(x0, y, color)
		}
	} else if y0 == y1 {
		// horizontal line
		xmin := min(x0, x1)
		xmax := max(x0, x1)
		for x := xmin; x <= xmax; x++ {
			img.SetNRGBA(x, y0, color)
		}
	} else {
		// sloped line
	}

}

func main() {
	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(100, 100, black)

	line(13, 20, 13, 20, img, white)
	line(13, 20, 13, 40, img, white)
	line(13, 20, 80, 20, img, white)
	line(13, 20, 80, 40, img, white)

	img = imaging.FlipV(img)
	err := imaging.Save(img, "output.png")
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
