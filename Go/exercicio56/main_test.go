package main

import "testing"

func Test_Multiplicattion(t *testing.T) {
	resAqr := Multiplicattion(5, 5)
	resAwt := 15

	if resAqr != resAwt {
		t.Errorf("The operation 3 x 5 failed: expected %d, aquired %d", resAwt, resAqr)
	}
}
