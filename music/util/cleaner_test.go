// test driver for util package
package util

import (
	"fmt"
	"testing"
)

type testVal struct {
	in           string
	out          string
	expectChange bool
}

func TestCleanUni(t *testing.T) {
	//return
	testVals := []testVal{
		{"Déjà Vu", "Deja Vu", true},
		{"Deja Vu", "Deja Vu", false},
		{"Déjà Vu", "Deja Vu", true},
		{"Eeye BÃ©", "Eeye Be", true},
		{"El salÃ³n MÃ©xico", "El Salon Mexico", true},
		{"Antonín Dvořák", "Antonin Dvorak", true},
		{"AntonÃ­n DvoÅ™Ã¡k", "Antonin Dvorak", true},
	}
	var done bool
	for _, v := range testVals {
		rval := CleanUni(v.in, &done)
		fmt.Printf("rvl: %s d:%t v0: %s  V1: %s %t\n", rval, done, v.in, v.out, v.expectChange)
		if rval != v.out {
			t.Errorf("not equal r: %s != %s", rval, v.in)
		}
		if done != v.expectChange {
			t.Errorf("change flag %s not as expected  %t != %t", v.in, v.expectChange, done)
		}
	}
}
func TestRemovePunct(t *testing.T) {
	testVals := []testVal{
		{"lovin'", "lovin", true},
		{"don't", "dont", true},
		{"can't", "cant", true},
		{"I'm a' fixn' to go", "Im a fixn to go", true},
	}
	for _, v := range testVals {
		rval, done := RemovePunct(v.in)
		if rval != v.out {
			t.Errorf("not equal %s", v.in)
		}
		if done != v.expectChange {
			t.Errorf("change flag on %s not as expected  %t != %t", v.in, v.expectChange, done)
		}
	}
}
