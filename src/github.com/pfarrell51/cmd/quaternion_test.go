// utilities (math) for quaternions

package quaternion

import (
	"math"
	"testing"
)

var (
	q1  = Quaternion{1, 0, 0, 0}
	q2  = Quaternion{a: 10}
	q3  = Quaternion{11, 0, 0, 0}
	q4  = Quaternion{10, 0, 0, 0}
	qv1 = Quaternion{0, 1, 0, 0}
	qv2 = Quaternion{0, 0, 1, 1}
	qv3 = Quaternion{0, 1, 0, 1}
	qv4 = Quaternion{0, 2, 1, 2}
	qv5 = Quaternion{0, -2, -1, -2}
	qv6 = Quaternion{-1, -1, 1, 1}
	qa1 = Quaternion{4, -4, -4, 4}
	qa2 = Quaternion{0.5, -0.5, -0.5, 0.5}
	qa3 = Quaternion{0.0625, 0.0625, 0.0625, -0.0625}
	qa4 = Quaternion{0.24765262787484427, 0.2940044459739585, 0.3943046179925829, 0.8347175749221727}
	qa5 = Quaternion{-0.7904669075670613, 0.44891659738265544, -0.3627631346111533, 0.205033813803568}
	qa6 = Quaternion{math.Cos(math.Pi / 2), math.Sin(math.Pi/2) / math.Sqrt(3),
		math.Sin(math.Pi/2) / math.Sqrt(3), -math.Sin(math.Pi/2) / math.Sqrt(3)}
	m = [3][3]float64{[3]float64{-0.333333333, 0.666666667, -0.666666667},
		[3]float64{0.666666667, -0.333333333, -0.666666667},
		[3]float64{-0.666666667, -0.666666667, -0.333333333}}
)

func TestMagnitude(t *testing.T) {
	mag := q1.Magnitude()
	if mag < 0 {
		t.Fail()
	}
}

func TestNorm(t *testing.T) {
	rval := q1.Norm()
	if rval.Magnitude() != 1 {
		t.Error("identity norm")
	}
	rval = qa1.Norm()
	if rval != qa2 {
		t.Error("4's norm")
	}
}

func TestVectorMag(t *testing.T) {
	if qv4.Magnitude() != 3 {
		t.Error("vector mag")
	}
}

func TestMixedMag(t *testing.T) {
	if qa1.Magnitude() != 8 {
		t.Errorf("mixed mag %g", qa5.Magnitude())
	}
}
func TestConj(t *testing.T) {
	rval := q1.Conj()
	if rval != q1 {
		t.Error("conj1")
	}
	rval = qv5.Conj()
	if qv4 != qv5.Conj() {
		t.Error("conj2")
	}
}
func TestSum1(t *testing.T) {
	rval := Sum(q1, q2)
	if rval != q3 {
		t.Errorf("sum1")
	}
}
func TestProd1(t *testing.T) {
	rval := Prod(q1, q2)
	if rval != q4 {
		t.Error("prod1a")
	}
	rval = Prod(q1, q1)
	if rval != q1 {
		t.Error("prod1b")
	}
}
func TestVectorProd(t *testing.T) {
	rval := Prod(qv1, qv2, qv3)
	if rval != qv6 {
		t.Errorf("vector prod  %s != %s", rval, qv6)
	}
}
func TestInv(t *testing.T) {
	if qa1.Inv() != qa3 {
		t.Error("Inv")
	}
}

func TestEuler(t *testing.T) {
	phi, theta, psi := Euler(qa4)
	if math.Abs(phi-1.0) > 1e-6 ||
		math.Abs(theta+0.3) > 1e-6 ||
		math.Abs(psi-2.4) > 1e-6 {
		t.Error("Euler")
	}
}
func TestFromEuler(t *testing.T) {
	q := FromEuler(-1.2, 0.4, 5.5)
	if math.Abs(q.a-qa5.a) > 1e-6 ||
		math.Abs(q.b-qa5.b) > 1e-6 ||
		math.Abs(q.c-qa5.c) > 1e-6 ||
		math.Abs(q.d-qa5.d) > 1e-6 {
		t.Errorf("FromEuler %s", q)
	}
}
func TestRotMat(t *testing.T) {
	mm := RotMat(qa6)
	for i, x := range mm {
		for j, y := range x {
			if math.Abs(m[i][j]-y) > 1e-6 {
				t.Error("Rot")
			}
		}
	}
}
