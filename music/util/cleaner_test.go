package util

import (
	"fmt"
	"testing"
)

func TestCleanUni(t *testing.T) {
	testNames := [][]string{
		{"Deja Vu", "Deja Vu"},
		{"El salÃ³n MÃ©xico", "El Salon Mexico"},
		{"Antonín Dvořák", "Antonin Dvorak"},
		{"AntonÃ­n DvoÅ™Ã¡k", "Antonin Dvorak"},
	}
	var done bool
	for i, v := range testNames {
		fmt.Printf("i: %d l(v) %d v: %s\n", i, len(v), v)
		rval := CleanUni(v[0], &done)
		fmt.Printf("rvl: %s v0: %s  V1: %s\n", rval, v[0], v[1])
		if rval != v[1] {
			t.Errorf("not equal ")
		}
	}
}
