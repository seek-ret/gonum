package goblas

import (
	"errors"
	"math"
)

const (
	debug = false
)

func newGeneral(r, c int) general {
	return general{
		data:   make([]float64, r*c),
		rows:   r,
		cols:   c,
		stride: c,
	}
}

type general struct {
	data       []float64
	rows, cols int
	stride     int
}

// adds element-wise into receiver. rows and columns must match
func (g general) add(h general) {
	if debug {
		if g.rows != h.rows {
			panic("row size mismatch")
		}
		if g.cols != h.cols {
			panic("col size mismatch")
		}
	}
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.data[i*g.stride+j] += h.at(i, j)
		}
	}
}

// at returns the value at the ith row and jth column. For speed reasons, the
// rows and columns are not bounds checked.
func (g general) at(i, j int) float64 {
	if debug {
		if i < 0 || i >= g.rows {
			panic("row out of bounds")
		}
		if j < 0 || j >= g.cols {
			panic("col out of bounds")
		}
	}
	return g.data[i*g.stride+j]
}

func (g general) check() error {
	if g.rows < 0 {
		return errors.New("general: rows < 0")
	}
	if g.cols < 0 {
		return errors.New("general: cols < 0")
	}
	if g.stride < 1 {
		return errors.New("general: stride < 1")
	}
	if g.stride < g.cols {
		return errors.New("general: illegal stride")
	}
	if (g.rows-1)*g.stride+g.cols > len(g.data) {
		return errors.New("general: insufficient length")
	}
	return nil
}

func (g general) clone() general {
	data := make([]float64, len(g.data))
	copy(data, g.data)
	return general{
		data:   data,
		rows:   g.rows,
		cols:   g.cols,
		stride: g.stride,
	}
}

// assumes they are the same size
func (g general) copy(h general) {
	if debug {
		if g.rows != h.rows {
			panic("row mismatch")
		}
		if g.cols != h.cols {
			panic("col mismatch")
		}
	}
	for k := 0; k < g.rows; k++ {
		copy(g.data[k*g.stride:(k+1)*g.stride], h.data[k*h.stride:(k+1)*h.stride])
	}
}

func (g general) equal(a general) bool {
	if g.rows != a.rows || g.cols != a.cols || g.stride != a.stride {
		return false
	}
	for i, v := range g.data {
		if a.data[i] != v {
			return false
		}
	}
	return true
}

/*
// print is to aid debugging. Commented out to avoid fmt import
func (g general) print() {
	fmt.Println("r = ", g.rows, "c = ", g.cols, "stride: ", g.stride)
	for i := 0; i < g.rows; i++ {
		fmt.Println(g.data[i*g.stride : (i+1)*g.stride])
	}

}
*/

func (g general) view(i, j, r, c int) general {
	if debug {
		if i < 0 || i+r > g.rows {
			panic("row out of bounds")
		}
		if j < 0 || j+c > g.cols {
			panic("col out of bounds")
		}
	}
	return general{
		data:   g.data[i*g.stride+j : (i+r-1)*g.stride+j+c],
		rows:   r,
		cols:   c,
		stride: g.stride,
	}
}

func (g general) equalWithinAbs(a general, tol float64) bool {
	if g.rows != a.rows || g.cols != a.cols || g.stride != a.stride {
		return false
	}
	for i, v := range g.data {
		if math.Abs(a.data[i]-v) > tol {
			return false
		}
	}
	return true
}
