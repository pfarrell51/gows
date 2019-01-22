// utilities (math) for quaternions

package quaternion

import (
	"fmt"
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
)

func TestMagnitude(t *testing.T) {
	mag := q1.Magnitude()
	fmt.Printf("mag: %g\n", mag)
	if mag < 0 {
		t.Fail()
	}

}

func TestNorm(t *testing.T) {
	rval := q1.Norm()
	if rval.Magnitude() != 1 {
		t.Fail()
	}
}
func TestConj(t *testing.T) {
	rval := q1.Conj()
	fmt.Printf("conj: %s\n", rval)
	if rval != q1 {
		t.Error("conj1")
	}
	rval = qv5.Conj()
	fmt.Printf("conj-b: %s\n", rval)
	if qv4 != qv5.Conj() {
		t.Error("conj2")
	}
}
func TestSum1(t *testing.T) {
	rval := Sum(q1, q2)
	fmt.Printf("in Sum1() %s \n", rval)
	if rval != q3 {
		t.Errorf("sum1")
	}
}
func TestProd1(t *testing.T) {
	rval := Prod(q1, q2)
	fmt.Printf("in Prod1() %s \n", rval)
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
