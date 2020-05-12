package main

import (
	"math"
	"testing"
)

func isPrimeEx(n float64) (factor float64, isPrime bool) {
	var i float64
	for i = 2; i <= math.Floor(math.Sqrt(n)); i++ {
		if math.Mod(n, i) == 0 {
			return i, false
		}
	}
	return 0, true
}

func TestIsPrime(t *testing.T) {
	factor, isPrime := isPrimeEx(21)

	if factor != 3 {
		t.Errorf("expected factor %v but received %v for argument %v", 3, factor, 21)
	}

	if isPrime != false {
		t.Errorf("expected isPrime %v but received %v for argument %v", false, isPrime, 21)
	}
}
