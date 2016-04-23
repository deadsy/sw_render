package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/deadsy/sw_render/wavefront"
	"github.com/disintegration/imaging"
)

const image_size = 1000

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
		fmt.Printf("%s: %s\n", objfile, err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", obj)

	ofs := [3]float32{-obj.Min_V(0), -obj.Min_V(1), -obj.Min_V(2)}
	fmt.Printf("ofs: %+v\n", ofs)

	x_range := obj.Range_V(0)
	y_range := obj.Range_V(1)
	z_range := obj.Range_V(2)

	fmt.Printf("range: %f %f %f\n", x_range, y_range, z_range)

	// scale by x
	s := image_size / x_range
	scale := [3]float32{s, s, s}
	fmt.Printf("%+v\n", scale)

	x_size := int(x_range*s) + 1
	y_size := int(y_range*s) + 1
	z_size := int(z_range*s) + 1

	fmt.Printf("%d %d %d\n", x_size, y_size, z_size)

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}

	// plotting x and y values, dropping z
	img := imaging.New(x_size, y_size, black)

	// iterate over the object faces
	for i := 0; i < obj.Len_F(); i++ {

		// get the vertices from the face
		v0 := obj.Get_V(i, 0)
		v1 := obj.Get_V(i, 1)
		v2 := obj.Get_V(i, 2)

		p0 := v0.Scale(&ofs, &scale)
		p1 := v1.Scale(&ofs, &scale)
		p2 := v2.Scale(&ofs, &scale)

		//fmt.Printf("%+v %+v %+v\n", p0, p1, p2)

		// p0 to p1
		line(V2{p0[0], p0[1]}, V2{p1[0], p1[1]}, img, white)
		// p1 to p2
		line(V2{p1[0], p1[1]}, V2{p2[0], p2[1]}, img, white)
		// p2 to p0
		line(V2{p2[0], p2[1]}, V2{p0[0], p0[1]}, img, white)

	}

	img = imaging.FlipV(img)
	err = imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
