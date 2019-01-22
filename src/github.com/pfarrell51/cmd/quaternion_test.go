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
	fmt.Printf("%g %g %g %g\n", rval.a, rval.b, rval.c, rval.d)
}
func TestSum1(t *testing.T) {
	rval := Sum(q1, q2)
	fmt.Printf("in Sum1() %g %g %g %g\n", rval.a, rval.b, rval.c, rval.d)
	if rval != q3 {
		t.Errorf("sum1")
	}

}
