package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/deadsy/sw_render/wavefront"
	"github.com/disintegration/imaging"
)

type V2 [2]int

func Equal(a, b V2) bool {
	return (a[0] == b[0]) && (a[1] == b[1])
}

func Equal_X(a, b V2) bool {
	return a[0] == b[0]
}

func Equal_Y(a, b V2) bool {
	return a[1] == b[1]
}

func Min_Y(a, b V2) V2 {
	if a[1] < b[1] {
		return a
	} else {
		return b
	}
}

func Max_Y(a, b V2) V2 {
	if a[1] < b[1] {
		return b
	} else {
		return a
	}
}

func Min_X(a, b V2) V2 {
	if a[0] < b[0] {
		return a
	} else {
		return b
	}
}

func Max_X(a, b V2) V2 {
	if a[0] < b[0] {
		return b
	} else {
		return a
	}
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func line(a, b V2, img *image.NRGBA, color color.NRGBA) {

	if Equal(a, b) {
		// point
		img.SetNRGBA(a[0], a[1], color)
	} else if Equal_X(a, b) {
		// vertical line
		p0 := Min_Y(a, b)
		p1 := Max_Y(a, b)
		for y := p0[1]; y <= p1[1]; y++ {
			img.SetNRGBA(a[0], y, color)
		}
	} else if Equal_Y(a, b) {
		// horizontal line
		p0 := Min_X(a, b)
		p1 := Max_X(a, b)
		for x := p0[0]; x <= p1[0]; x++ {
			img.SetNRGBA(x, a[1], color)
		}
	} else {
		// sloped line
		if abs(a[0]-b[0]) > abs(a[1]-b[1]) {
			// x is the long axis
			p0 := Min_X(a, b)
			p1 := Max_X(a, b)
			dx := p1[0] - p0[0]
			dy := p1[1] - p0[1]
			y := p0[1]
			err := 0
			d_err := 2 * dy
			for x := p0[0]; x <= p1[0]; x++ {
				img.SetNRGBA(x, y, color)
				err += d_err
				if err >= dx {
					y += 1
					err -= 2 * dx
				}
			}
		} else {
			// y is the long axis
			p0 := Min_Y(a, b)
			p1 := Max_Y(a, b)
			dx := p1[0] - p0[0]
			dy := p1[1] - p0[1]
			x := p0[0]
			err := 0
			d_err := 2 * dx
			for y := p0[1]; y <= p1[1]; y++ {
				img.SetNRGBA(x, y, color)
				err += d_err
				if err >= dy {
					x += 1
					err -= 2 * dy
				}
			}
		}
	}
}

func main() {

	//objfile := "../obj/gopher.obj"
	objfile := "../obj/african_head.obj"
	imgfile := "output.png"

	obj, err := wavefront.Read(objfile)
	if err != nil {
		fmt.Printf("unable to load %s, %s\n", objfile, err)
		os.Exit(1)
	}

	obj.Display()

	os.Exit(0)

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(100, 100, black)

	line(V2{13, 20}, V2{13, 20}, img, white)
	line(V2{13, 20}, V2{13, 40}, img, white)
	line(V2{13, 20}, V2{80, 20}, img, white)
	line(V2{13, 20}, V2{80, 40}, img, white)
	line(V2{13, 20}, V2{40, 80}, img, white)
	line(V2{0, 0}, V2{99, 1}, img, white)
	line(V2{0, 0}, V2{99, 5}, img, white)
	line(V2{0, 0}, V2{99, 99}, img, white)
	line(V2{0, 0}, V2{1, 99}, img, white)

	img = imaging.FlipV(img)
	err = imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
