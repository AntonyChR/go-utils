package utils

import (
	"math"
	rand "math/rand"
)

type Vec []float64

func (v *Vec) Equals(b Vec) bool {
	if len(*v) != len(b) {
		return false
	}

	for i, v := range b {
		if v != b[i] {
			return false
		}
	}
	return true
}

func (v *Vec) Len() int {
	return len(*v)
}

func (v *Vec) Unit() Vec {
	l := v.Norm()
	size := v.Len()

	if l == 0 {
		return Zero(size)
	}

	c := make([]float64, size)

	for i, e := range *v {
		c[i] = e / l
	}
	return c
}

func (v *Vec) Norm() float64 {
	acc := 0.0
	for _, el := range *v {
		acc += el * el
	}
	return math.Sqrt(acc)
}

func (v *Vec) Add(value interface{}) Vec {
	if u, ok := value.(float64); ok {
		c := make([]float64, v.Len())
		for i, n := range *v {
			c[i] = n + u
		}
		return c
	}
	if u, ok := value.(Vec); ok {
		AssertEq(len(u), v.Len(), "Vectors must have the same size")
		c := make([]float64, v.Len())
		for i := range v.Len() {
			c[i] = u[i] + (*v)[i]
		}
		return c
	}
	panic("The parameter 'value' must be float64 or an instance of Vec")
}

func (v *Vec) Sub(value interface{}) Vec {
	if u, ok := value.(float64); ok {
		c := make([]float64, v.Len())
		for i, n := range *v {
			c[i] = n - u
		}
		return c
	}
	if u, ok := value.(Vec); ok {
		AssertEq(len(u), v.Len(), "Vectors must have the same size")
		c := make([]float64, v.Len())
		for i := range v.Len() {
			c[i] = (*v)[i] - u[i]
		}
		return c
	}
	panic("The parameter 'value' must be float64 or an instance of Vec")
}

func (v *Vec) Scale(n float64) Vec {
	c := make([]float64, v.Len())

	for i := range *v {
		c[i] = (*v)[i] * n
	}
	return c
}

func (v *Vec) Dot(u Vec) float64 {
	AssertEq(len(u), v.Len(), "Vectors must have the same size")
	acc := 0.0
	for i := range *v {
		acc += (*v)[i] * u[i]
	}
	return acc
}

func (v *Vec) DistanceTo(u Vec) float64 {
	t := v.Sub(u)
	return t.Norm()
}

func (v *Vec) Fill(value float64) {
	for i := range v.Len() {
		(*v)[i] = value
	}
}

func (v *Vec) Inv() Vec {
	l := v.Len()
	c := make([]float64, l)
	for i, e := range *v {
		c[l-i-1] = e
	}
	return c
}

func (v *Vec) Resize(size int) Vec {
	if size < v.Len() {
		return (*v)[:size]
	}
	c := make([]float64, size)
	copy(c, *v)
	return c
}

func (v *Vec) Concat(u Vec) Vec {
	return append((*v), u...)
}
func (v *Vec) Map(cb func(float64) float64) Vec {
	c := make([]float64, v.Len())
	for i, e := range *v {
		c[i] = cb(e)
	}
	return c
}

//
// methods for 2d *vetors
//

const ERROR_2d_METHOD = "This method is only avalaible for 2d *vetors"

func (v *Vec) ProjectOn(u Vec) Vec {
	AssertEq(v.Len(), 2, ERROR_2d_METHOD)
	unit := u.Unit()
	return unit.Scale(v.Dot(u) / u.Norm())
}

func (v *Vec) Rot(angle float64) Vec {
	AssertEq(v.Len(), 2, ERROR_2d_METHOD)
	x := (*v)[0]*math.Cos(angle) - (*v)[1]*math.Sin(angle)
	y := (*v)[0]*math.Sin(angle) + (*v)[1]*math.Cos(angle)
	return Vec{x, y}
}

func (v *Vec) IsParallelTo(u Vec) bool {
	AssertEq(v.Len(), 2, ERROR_2d_METHOD)
	AssertEq(u.Len(), 2, ERROR_2d_METHOD)

	return math.Abs((*v)[0]*u[1]-(*v)[1]*u[0]) < 1e-9
}

func Vec2FromAngle(angle float64) Vec {
	return Vec{math.Cos(angle), math.Sin(angle)}
}

//
// methods for 3d vectors
//

const ERROR_3D_METHOD = "This method is only available for 3d vectors"

func (v *Vec) Prod3d(u Vec) Vec {
	AssertEq(v.Len(), 3, ERROR_3D_METHOD)
	AssertEq(u.Len(), 3, ERROR_3D_METHOD)

	return Vec{
		(*v)[1]*u[2] - (*v)[2]*u[1],
		(*v)[2]*u[0] - (*v)[0]*u[2],
		(*v)[0]*u[1] - (*v)[1]*u[0],
	}
}

func Zero(size int) Vec {
	return make([]float64, size)
}

func One(size int) Vec {
	return Filled(size, 1)
}

func Filled(size int, value float64) Vec {
	c := make([]float64, size)
	for i := range size {
		c[i] = value
	}
	return c
}

func Rand(size int) Vec {
	c := make([]float64, size)
	for i := range size {
		c[i] = rand.Float64()
	}
	return c
}
