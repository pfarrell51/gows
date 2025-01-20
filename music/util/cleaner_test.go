// test driver for util package
package util

import (
	//"fmt"
	"testing"
)

type testVal struct {
	in           string
	out          string
	expectChange bool
}

func TestCleanUni(t *testing.T) {
	testVals := []testVal{
		{"bànãnà", "banana", true},
		{"bànãnà dànãnà", "banana danana", true},
		{"Déjà Vu", "Deja Vu", true},
		{"Deja Vu", "Deja Vu", false},
		{"Déjà Vú", "Deja Vu", true},
		{"Eeye BÃ©", "Eeye Be", true},
		{"El SalÃ³n MÃ©xico", "El Salon Mexico", true},
		{"Antonín Dvořák", "Antonin Dvorak", true},
		//	{"AntonÃ­n DvoÅ™Ã¡k", "Antonin Dvorak", true},
	}
	for _, v := range testVals {
		var done bool
		rval := CleanUni(v.in, &done)
		if rval != v.out {
			t.Errorf("not equal r: %s != %s for %s", rval, v.out, v.in)
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
			t.Errorf("not equal in %s rval %s", v.in, rval)
		}
		if done != v.expectChange {
			t.Errorf("change flag on %s not as expected  %t != %t", v.in, v.expectChange, done)
		}
	}
}
