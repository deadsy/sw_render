package vec

type V2i [2]int

// return true if the vectors are equal
func (a V2i) Equal(b V2i) bool {
	return (a[0] == b[0]) && (a[1] == b[1])
}

// Return a - b
func (a V2i) Sub(b V2i) V2i {
	return V2i{
		a[0] - b[0],
		a[1] - b[1],
	}
}
