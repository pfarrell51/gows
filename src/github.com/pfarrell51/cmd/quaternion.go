// utilities (math) for Quaternions

package quaternion

import (
	"fmt"
	"math"
)

type Quaternion struct {
	a, b, c, d float64
}

func (q Quaternion) String() string {
	return fmt.Sprintf("[%g, %gi, %gj, %gk]",
		q.a, q.b, q.c, q.d)
}

// The magnitude or length of a Quaternion a + i b + j c+k d is sqrt ( a2 + b2+c2+ d2)
// magnitude(a + i b + c j + d k) = sqrt ( a**2 + b** 2+c**2+ d**2)

func (q Quaternion) Magnitude() float64 {
	var rval float64
	rval = math.Sqrt(q.a*q.a + q.b*q.b + q.c*q.c + q.d*q.d)
	return rval
}

// Norm = || q || = sqrt( q * conj(q)) = sqrt ( a**2 + b**2+c**2+ d**2)
// returns a normalized Quaternion with a length of 1
// divide each of a,b,c and d by the length of the vector, this will make || q || = 1.
func (q Quaternion) Norm() Quaternion {
	mag := q.Magnitude()
	rval := Quaternion{q.a / mag, q.b / mag, q.c / mag, q.d / mag}
	return rval
}

// returns conjugate of the aregument (same real, reverse since of i,j,k)
func (q Quaternion) Conj() Quaternion {
	rval := Quaternion{}
	rval.a = q.a
	if rval.b != -0.0 {
		rval.b = -q.b
	}
	if rval.c != -0.0 {
		rval.c = -q.c
	}
	if rval.d != -0.0 {
		rval.d = -q.d
	}
	return rval
}

// non-commutatve product of any number of q's
func Prod(ain ...Quaternion) Quaternion {
	rval := Quaternion{1, 0, 0, 0}
	var w, x, y, z float64
	for _, q := range ain {
		fmt.Printf("pl: %s\n", q)
		w = rval.a*q.a - rval.b*q.b - rval.c*q.c - rval.d*q.d
		x = rval.a*q.b + rval.b*q.a + rval.c*q.d - rval.d*q.c
		y = rval.a*q.c + rval.c*q.a + rval.d*q.b - rval.a*q.d
		z = rval.a*q.d + rval.d*q.a + rval.b*q.c - rval.c*q.b
		rval = Quaternion{w, x, y, z}
	}
	return rval
}

// som of any number of q's
func Sum(ain ...Quaternion) Quaternion {
	rval := Quaternion{}
	for _, q := range ain {
		rval.a += q.a
		rval.b += q.b
		rval.c += q.c
		rval.d += q.d
	}
	return rval

}
