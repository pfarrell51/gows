// utilities (math) for latlong

package main

import (
	"fmt"
	"testing"
)

var (
	pat      = Latlong{40.064587, -75.37118}
	doug     = Latlong{35.151035, -81.0029}
	pittrace = Latlong{40.844715, -80.32963}
)

func TestDistance(t *testing.T) {
	rval := pat.distance(doug)
	fmt.Printf("to doug %g\n ", rval)
	if rval < 5900.0 {
		t.Fail()
	}
	rval = pat.distance(pittrace)
	fmt.Printf("to Pitt %g\n", rval)
}
