// utilities (math) for quaternions

package quaternion

import (
	"fmt"
	"testing"
)

var (
	q1 = Quaternion{1, 0, 0, 0}
	q2 = Quaternion{a: 10}
	q3 = Quaternion{11, 0, 0, 0}
	q4 = Quaternion{10, 0, 0, 0}
)

func TestNegZero(t *testing.T) {
	var foo float64
	foo = 0
	fmt.Printf("pos %18.8e\n", foo)
	foo = -foo
	fmt.Printf("neg %18.8e\n", foo)
}
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
