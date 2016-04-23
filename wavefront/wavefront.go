package wavefront

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type file_error struct {
	msg string
}

func (e *file_error) Error() string {
	return e.msg
}

// geometric vertex
type V_elem struct {
	x [3]float32
	w float32
}

// texture coordinate
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
	vt int // index to texture coordinate
	vn int // index for vertex normal
}

// Object contains a name and a list of groups
type Object struct {
	v_list  []*V_elem
	vt_list []*VT_elem
	vn_list []*VN_elem
	f_list  []*[3]F_elem
}

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

// validate the f element indices
func (o *Object) Check_F(f *F_elem) error {
	if f.v > len(o.v_list) {
		return &file_error{"v index out of range"}
	}
	if f.vt > len(o.vt_list) {
		return &file_error{"vt index out of range"}
	}
	if f.vn > len(o.vn_list) {
		return &file_error{"vn index out of range"}
	}
	return nil
}

func (o *Object) Display() {
	fmt.Printf("geometric vertices %d\n", len(o.v_list))
	fmt.Printf("texture coordinates %d\n", len(o.vt_list))
	fmt.Printf("vertex normals %d\n", len(o.vn_list))
	fmt.Printf("faces %d\n", len(o.f_list))
}

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
		return fmt.Errorf(msg+" at %s:%d: %s", filename, line_number, line)
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
			copy(v.x[:], x[0:2])
			object.Add_V(&v)

		case "vt":
			// texture coordinate
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
			copy(vt.x[:], x[0:1])
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

		case "g":
			// group name

		case "s":
			// smooth shading

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
					f[i].v = x
					// index to texture coordinate (could be empty)
					if len(indices) >= 2 {
						if len(indices[1]) > 0 {
							x, err := strconv.Atoi(indices[1])
							if err != nil {
								return nil, fail("bad face texture coordinate")
							}
							f[i].vt = x
						}
						// index to vertex normal
						if len(indices) >= 3 {
							x, err := strconv.Atoi(indices[2])
							if err != nil {
								return nil, fail("bad face vertex normal")
							}
							f[i].vn = x
						}
					}
				} else {
					return nil, fail("bad face indices")
				}
				err := object.Check_F(&f[i])
				if err != nil {
					return nil, fail(err.Error())
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
