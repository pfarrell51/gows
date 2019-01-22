// utilities (math) for Quaternions

package quaternion

import (
	"math"
)

type Quaternion struct {
	a, b, c, d float64
}

// The magnitude or length of a Quaternion a + i b + j c+k d is sqrt ( a2 + b2+c2+ d2)
// magnitude(a + i b + c j + d k) = sqrt ( a**2 + b** 2+c**2+ d**2)

func (q Quaternion) Magnitude() float64 {
	var rval float64
	rval = math.Sqrt(q.a*q.a + q.b*q.b + q.c*q.c + q.d*q.d)
	return rval
}

// Norm = || q || = sqrt( q * conj(q)) = sqrt ( a**2 + b**2+c**2+ d**2)
// To Normalise a Quaternion divide each of a,b,c and d by the above value, this will make || q || = 1.
func (q Quaternion) Norm() Quaternion {
	mag := q.Magnitude()
	rval := Quaternion{}
	rval.a = q.a / mag
	rval.b = q.b / mag
	rval.c = q.c / mag
	rval.d = q.d / mag
	return rval
}

// returns conjugate of the aregument (same real, reverse since of i,j,k)
func (q Quaternion) Conj() Quaternion {
	rval := Quaternion{}
	rval.a = q.a
	rval.b = -q.b
	rval.c = -q.c
	rval.d = q.d
	return rval
}
func Sum(a1, a2 Quaternion) Quaternion {
	rval := Quaternion{}
	rval.a = a1.a + a2.a
	rval.b = a1.b + a2.b
	rval.a = a1.c + a2.c
	rval.a = a1.d + a2.d
	return rval

}
