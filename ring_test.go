package abp

import (
	"testing"
)

func TestNegativeSize(t *testing.T) {
	r := newRing(-1)
	if r != nil {
		t.Errorf("Ring should be nil but was %v", r)
	}
}

func TestMove(t *testing.T) {
	r := newRing(3)
	for i := 1; i < 4; i++ {
		r.value = i
		r = r.Move(1)
	}
	if r.value != 1 {
		t.Errorf("Ring values should be 1 but was %d", r.value)
	}

	r = r.Move(-1)
	if r.value != 3 {
		t.Errorf("Ring values should be 3 but was %d", r.value)
	}
}
