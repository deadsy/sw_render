# Software Rendering

Learning more about polygon rendering
https://github.com/ssloy/tinyrenderer/wiki
... in golang

## Lesson 0
Using my language of choice (Go) output an image with a pixel set.

![lesson0_image0](https://github.com/deadsy/sw_render/blob/master/lesson0/pics/test.jpeg "lesson0_image0")

## Lesson 1
### First Head Image
Apparently lines with negative slope are a challenge.

![lesson1_image0](https://github.com/deadsy/sw_render/blob/master/lesson1/pics/1.png "lesson1_image0")

After I covered all of the cases for line drawing properly.

![lesson1_image1](https://github.com/deadsy/sw_render/blob/master/lesson1/pics/2.png "lesson1_image1")

And the wireframe for the go gopher - axis choices put it on it's side.

![lesson1_image2](https://github.com/deadsy/sw_render/blob/master/lesson1/pics/3.png "lesson1_image2")

## Lesson 2

Wherein we color in triangles...

Testing triangle rasterization based on Bresenham's line drawing algorithm.

![lesson2_image1](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/1.png "lesson2_image1")

Long thin triangles are a challenge. I'm incrementing in the y direction, but because dx > dy for both
of the lines I end up with some gaps where the x step is >= 2. Not entirely sure what the right thing
to do is- but clearly having a contiguous region rendered would be nice. 

e.g. 	triangle(vec.V2i{0, 0}, vec.V2i{5, 0}, vec.V2i{300, 100}, img, white)

![lesson2_image2](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/2.png "lesson2_image2")

Random triangles - quite a few with with the aforementioned dots on the vertices problem.

![lesson2_image3](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/3.png "lesson2_image3")

Random triangles for the head mesh. The isolated pixels problem is not immediately apparent here. I 
suppose this is because the head has been created with triangles that are more well behaved.

![lesson2_image4](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/4.png "lesson2_image4")

Random thoughts on triangle rendering:

With the Bresenham style rendering I am doing I end up with (by design) a pixels at each vertex.
Clearly a triangle could be so small that it doesn't even deserve a single pixel to represent it.

Another idea: Bounding box, followed by scanning the horizontal pixels. If the pixel centroid is
within the triangle then turn it on. This should be a contiguous region- albeit still problematic
for long thin triangles, where two triangles of similar shape might have very different renderings.

This time with random coloring of the triangles with normals facing forwards.

![lesson2_image5](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/5.png "lesson2_image5")

And with grey scale shading based on the dot product of the light and the face normal vector.

![lesson2_image6](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/6.png "lesson2_image6")

Flashlight under the chin...

![lesson2_image7](https://github.com/deadsy/sw_render/blob/master/lesson2/pics/7.png "lesson2_image7")



