package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"

	"github.com/deadsy/sw_render/vec"
	"github.com/deadsy/sw_render/wavefront"
	"github.com/disintegration/imaging"
)

func Random_Color() color.NRGBA {
	return color.NRGBA{
		uint8(256 * rand.Float32()),
		uint8(256 * rand.Float32()),
		uint8(256 * rand.Float32()),
		255,
	}
}

func Grey_Scale(level float32) color.NRGBA {
	return color.NRGBA{
		uint8(255 * level),
		uint8(255 * level),
		uint8(255 * level),
		255,
	}
}

func random_triangles(k vec.V2i, img *image.NRGBA) {
	for i := 0; i < 200; i++ {
		a := k.Rand()
		b := a.Rand_Delta(100)
		c := a.Rand_Delta(100)
		triangle(a, b, c, img, Random_Color())
	}
}

const pixels_x = 750
const pixels_ofs = 5

// object to image mapping
func Obj2Img(v, ofs vec.V3, scale float32) vec.V2i {
	p := v.Sum(ofs).Scale(scale)
	return vec.V2i{int(p[0]), int(p[1])}
}

func main() {
	test_barycentric()
}

func main2() {

	imgfile := "output.png"
	objfile := "../obj/african_head.obj"

	obj, err := wavefront.Read(objfile)
	if err != nil {
		fmt.Printf("%s: %s\n", objfile, err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", obj)

	obj_ofs := obj.Offset()
	fmt.Printf("obj_ofs: %+v\n", obj_ofs)

	obj_range := obj.Range()
	fmt.Printf("obj_range: %+v\n", obj_range)

	// we want the image to be pixels_x wide on the x-axis
	// work out the scaling factor
	scale := (pixels_x - pixels_ofs) / obj_range[0]

	// work out the image size
	img_size := obj_range.Scale(scale)
	img_size = img_size.Sum(vec.V3{pixels_ofs, pixels_ofs, pixels_ofs})
	fmt.Printf("img_size: %+v\n", img_size)

	black := color.NRGBA{0, 0, 0, 255}
	img := imaging.New(int(img_size[0]), int(img_size[1]), black)
	light := vec.V3{0, 0, -1}.Normalize()

	// iterate over the object faces
	for i := 0; i < obj.Len_F(); i++ {

		// get the vertices from the face
		v0 := obj.Get_V(i, 0).ToV3()
		v1 := obj.Get_V(i, 1).ToV3()
		v2 := obj.Get_V(i, 2).ToV3()

		normal := v2.Sub(v0).Cross(v1.Sub(v0)).Normalize()
		shading := light.Dot(normal)

		if shading > 0 {
			p0 := Obj2Img(v0, obj_ofs, scale)
			p1 := Obj2Img(v1, obj_ofs, scale)
			p2 := Obj2Img(v2, obj_ofs, scale)
			triangle(p0, p1, p2, img, Grey_Scale(shading))
		}

	}

	//random_triangles(k, img)

	img = imaging.FlipV(img)
	err = imaging.Save(img, imgfile)
	if err != nil {
		fmt.Printf("unable to save %s, %s\n", imgfile, err)
		os.Exit(1)
	}

	os.Exit(0)
}
