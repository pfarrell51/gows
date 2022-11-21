// utilities (math) for Quaternions
// developed by Irish mathematician William Rowan Hamilton in 1843
// code inspired by work of Eric Westphal

package quaternion

import (
	"fmt"
	"math"
)

type Quaternion struct {
	a, // scalar
	b, // i
	c, // j
	d float64 // k
}

func (q Quaternion) String() string {
	return fmt.Sprintf("[%g, %gi, %gj, %gk]",
		q.a, q.b, q.c, q.d)
}

// The magnitude or length of a Quaternion a + i b + j c+k d is sqrt ( a2 + b2+c2+ d2)
// magnitude(a + i b + c j + d k) = sqrt ( a**2 + b** 2+c**2+ d**2)

func (q Quaternion) Magnitude() float64 {
	rval := math.Sqrt(q.L2Magnitude())
	return rval
}

// L2Magnitude returns the "L2-Norm" of a Quaternion (W,X,Y,Z) -> W*W+X*X+Y*Y+Z*Z
func (q Quaternion) L2Magnitude() float64 {
	rval := q.a*q.a + q.b*q.b + q.c*q.c + q.d*q.d
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

// returns conjugate of the aregument (same real, reverse sense of i,j,k)
func (q Quaternion) Conj() Quaternion {
	rval := Quaternion{q.a, q.b, q.c, q.d}
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

// Inv returns the Quaternion conjugate rescaled so that Q Q* = 1
func (qin Quaternion) Inv() Quaternion {
	k2 := qin.L2Magnitude()
	q := qin.Conj()
	return Quaternion{q.a / k2, q.b / k2, q.c / k2, q.d / k2}
}

// non-commutatve product of any number of q's
func Prod(ain ...Quaternion) Quaternion {
	rval := Quaternion{1, 0, 0, 0}
	var w, x, y, z float64
	for _, q := range ain {
		w = rval.a*q.a - rval.b*q.b - rval.c*q.c - rval.d*q.d
		x = rval.a*q.b + rval.b*q.a + rval.c*q.d - rval.d*q.c
		y = rval.a*q.c + rval.c*q.a + rval.d*q.b - rval.b*q.d
		z = rval.a*q.d + rval.d*q.a + rval.b*q.c - rval.c*q.b
		rval = Quaternion{w, x, y, z}
	}
	return rval
}

// sum of any number of q's
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

// Euler returns the Euler angles phi, theta, psi corresponding to a Quaternion
func Euler(q Quaternion) (float64, float64, float64) {
	r := q.Norm()
	phi := math.Atan2(2*(r.a*r.b+r.c*r.d), 1-2*(r.b*r.b+r.c*r.c))
	theta := math.Asin(2 * (r.a*r.c - r.d*r.b))
	psi := math.Atan2(2*(r.b*r.c+r.a*r.d), 1-2*(r.c*r.c+r.d*r.d))
	return phi, theta, psi
}

// FromEuler returns a Quaternion corresponding to Euler angles phi, theta, psi
func FromEuler(phi, theta, psi float64) Quaternion {
	q := Quaternion{}
	q.a = math.Cos(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) +
		math.Sin(phi/2)*math.Sin(theta/2)*math.Sin(psi/2)
	q.b = math.Sin(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) -
		math.Cos(phi/2)*math.Sin(theta/2)*math.Sin(psi/2)
	q.c = math.Cos(phi/2)*math.Sin(theta/2)*math.Cos(psi/2) +
		math.Sin(phi/2)*math.Cos(theta/2)*math.Sin(psi/2)
	q.d = math.Cos(phi/2)*math.Cos(theta/2)*math.Sin(psi/2) -
		math.Sin(phi/2)*math.Sin(theta/2)*math.Cos(psi/2)
	return q
}

// RotMat returns the rotation matrix (as float array) corresponding to a Quaternion
func RotMat(qin Quaternion) [3][3]float64 {
	q := qin.Norm()
	m := [3][3]float64{}
	m[0][0] = 1 - 2*(q.c*q.c+q.d*q.d)
	m[0][1] = 2 * (q.b*q.c - q.a*q.d)
	m[0][2] = 2 * (q.a*q.c + q.b*q.d)

	m[1][1] = 1 - 2*(q.d*q.d+q.b*q.b)
	m[1][2] = 2 * (q.c*q.d - q.a*q.b)
	m[1][0] = 2 * (q.a*q.d + q.c*q.b)

	m[2][2] = 1 - 2*(q.b*q.b+q.c*q.c)
	m[2][0] = 2 * (q.d*q.b - q.a*q.c)
	m[2][1] = 2 * (q.a*q.b + q.d*q.c)
	return m
}
