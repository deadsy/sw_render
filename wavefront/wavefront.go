//-----------------------------------------------------------------------------

package wavefront

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deadsy/sw_render/vec"
)

//-----------------------------------------------------------------------------

// geometric vertex
type V_elem struct {
	x [3]float32
	w float32
}

// texture vertex
type VT_elem struct {
	x [2]float32
	w float32
}

// vertex normal
type VN_elem struct {
	x [3]float32
}

type F_elem struct {
	v  int // index for geometric vertex
	vt int // index to texture vertex
	vn int // index for vertex normal
}

// Object contains a name and a list of groups
type Object struct {
	v_list  []*V_elem
	vt_list []*VT_elem
	vn_list []*VN_elem
	f_list  []*[3]F_elem
}

//-----------------------------------------------------------------------------
// operations on geometric vertices

// offset and scale
func (v *V_elem) Scale(ofs, scale *[3]float32) *[3]int {
	return &[3]int{
		int((v.x[0] + ofs[0]) * scale[0]),
		int((v.x[1] + ofs[1]) * scale[1]),
		int((v.x[2] + ofs[2]) * scale[2]),
	}
}

// convert a vertex to a V3
func (v *V_elem) ToV3() vec.V3f {
	return vec.V3f{v.x[0], v.x[1], v.x[2]}
}

//-----------------------------------------------------------------------------
// operations on objects

func (o *Object) Add_V(v *V_elem) {
	o.v_list = append(o.v_list, v)
}

func (o *Object) Add_VT(vt *VT_elem) {
	o.vt_list = append(o.vt_list, vt)
}

func (o *Object) Add_VN(vn *VN_elem) {
	o.vn_list = append(o.vn_list, vn)
}

func (o *Object) Add_F(f *[3]F_elem) {
	o.f_list = append(o.f_list, f)
}

func (o *Object) Len_V() int {
	return len(o.v_list)
}

func (o *Object) Len_VT() int {
	return len(o.vt_list)
}

func (o *Object) Len_VN() int {
	return len(o.vn_list)
}

func (o *Object) Len_F() int {
	return len(o.f_list)
}

// return the j-th vertex from the i-th face
func (o *Object) Get_V(i, j int) *V_elem {
	return o.v_list[o.f_list[i][j].v-1]
}

func (o *Object) String() string {
	var s []string
	s = append(s, fmt.Sprintf("geometric vertices %d", len(o.v_list)))
	s = append(s, fmt.Sprintf("texture vertices %d", len(o.vt_list)))
	s = append(s, fmt.Sprintf("vertex normals %d", len(o.vn_list)))
	s = append(s, fmt.Sprintf("faces %d", len(o.f_list)))
	s = append(s, fmt.Sprintf("bounds-x %f %f", o.Min_V(0), o.Max_V(0)))
	s = append(s, fmt.Sprintf("bounds-y %f %f", o.Min_V(1), o.Max_V(1)))
	s = append(s, fmt.Sprintf("bounds-z %f %f", o.Min_V(2), o.Max_V(2)))
	return strings.Join(s, "\n")
}

func (o *Object) Offset() vec.V3f {
	return vec.V3f{-o.Min_V(0), -o.Min_V(1), -o.Min_V(2)}
}

func (o *Object) Range() vec.V3f {
	return vec.V3f{o.Range_V(0), o.Range_V(1), o.Range_V(2)}
}

func (o *Object) Range_V(j int) float32 {
	return o.Max_V(j) - o.Min_V(j)
}

// return the value of the v element with the largest j-th index value
func (o *Object) Max_V(j int) float32 {
	x := float32(0)
	if len(o.v_list) != 0 {
		// initial value
		x = o.v_list[0].x[j]
		for i := 0; i < len(o.v_list); i++ {
			if o.v_list[i].x[j] > x {
				x = o.v_list[i].x[j]
			}
		}
	}
	return x
}

// return the value of the v element with the smallest j-th index value
func (o *Object) Min_V(j int) float32 {
	x := float32(0)
	if len(o.v_list) != 0 {
		// initial value
		x = o.v_list[0].x[j]
		for i := 0; i < len(o.v_list); i++ {
			if o.v_list[i].x[j] < x {
				x = o.v_list[i].x[j]
			}
		}
	}
	return x
}

//-----------------------------------------------------------------------------

func Read(filename string) (*Object, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	line_number := 0
	line := ""
	scanner := bufio.NewScanner(file)

	fail := func(msg string) error {
		return fmt.Errorf(msg+" at line %d", line_number)
	}

	var object Object

	for scanner.Scan() {
		line_number++
		line = scanner.Text()

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, " ") {
			continue
		}

		fields := strings.Fields(line)
		n_fields := len(fields)

		if n_fields == 0 {
			continue
		}

		n_fields -= 1

		switch fields[0] {

		case "o":
			// object name
			// TODO

		case "g":
			// group name
			// TODO

		case "s":
			// smoothing object
			// TODO

		case "mtllib":
			// material library
			// TODO

		case "usemtl":
			// use material
			// TODO

		case "l":
			// line
			// TODO

		case "v":
			// geometric vertex
			if n_fields < 3 {
				return nil, fail("v: not enough fields")
			}
			if n_fields > 4 {
				return nil, fail("v: too many fields")
			}
			var x [4]float32
			x[3] = 1
			for i := 0; i < n_fields; i++ {
				f, err := strconv.ParseFloat(fields[i+1], 32)
				if err != nil {
					return nil, fail("cannot parse float")
				}
				x[i] = float32(f)
			}
			v := V_elem{w: x[3]}
			copy(v.x[:], x[0:3])
			object.Add_V(&v)

		case "vt":
			// texture vertex
			if n_fields < 2 {
				return nil, fail("vt: not enough fields")
			}
			if n_fields < 3 {
				return nil, fail("vt: too many fields")
			}
			var x [3]float32
			x[2] = 0
			for i := 0; i < n_fields; i++ {
				f, err := strconv.ParseFloat(fields[i+1], 32)
				if err != nil {
					return nil, fail("cannot parse float")
				}
				x[i] = float32(f)
			}
			vt := VT_elem{w: x[2]}
			copy(vt.x[:], x[0:2])
			object.Add_VT(&vt)

		case "vn":
			// vertex normal
			if n_fields != 3 {
				return nil, fail("vn: wrong number of fields")
			}
			var vn VN_elem
			for i := 0; i < 3; i++ {
				f, err := strconv.ParseFloat(fields[i+1], 32)
				if err != nil {
					return nil, fail("cannot parse float")
				}
				vn.x[i] = float32(f)
			}
			object.Add_VN(&vn)

		case "f":
			// face
			if n_fields != 3 {
				return nil, fail("f: wrong number of fields")
			}
			var f [3]F_elem
			for i := 0; i < 3; i++ {
				indices := strings.Split(fields[i+1], "/")
				if len(indices) >= 1 {
					// index to geometric vertex
					x, err := strconv.Atoi(indices[0])
					if err != nil {
						return nil, fail("bad face geometric vertex index")
					}
					if x > object.Len_V() {
						return nil, fail("v index out of range")
					}
					f[i].v = x
					// index to texture vertex (could be empty)
					if len(indices) >= 2 {
						if len(indices[1]) > 0 {
							x, err := strconv.Atoi(indices[1])
							if err != nil {
								return nil, fail("bad face texture vertex")
							}
							if x > object.Len_VT() {
								return nil, fail("vt index out of range")
							}
							f[i].vt = x
						}
						// index to vertex normal
						if len(indices) >= 3 {
							x, err := strconv.Atoi(indices[2])
							if err != nil {
								return nil, fail("bad face vertex normal")
							}
							if x > object.Len_VN() {
								return nil, fail("vn index out of range")
							}
							f[i].vn = x
						}
					}
				} else {
					return nil, fail("bad face indices")
				}
			}
			object.Add_F(&f)

		default:
			return nil, fail("unrecognized element")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &object, nil
}

//-----------------------------------------------------------------------------
