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
		r = r.move(1)
	}
	if r.value != 1 {
		t.Errorf("Ring values should be 1 but was %d", r.value)
	}

	r = r.move(-1)
	if r.value != 3 {
		t.Errorf("Ring values should be 3 but was %d", r.value)
	}
}

func benchmarkMove(i int, b *testing.B) {
	r := newRing(i)
	for n := 0; n < b.N; n++ {
		r.move(1)
	}
}

func BenchmarkMove1000(b *testing.B) {
	benchmarkMove(1000, b)
}

func BenchmarkMove10000(b *testing.B) {
	benchmarkMove(10000, b)
}

func benchmarkMoveFull(i int, b *testing.B) {
	r := newRing(i)
	for n := 0; n < b.N; n++ {
		r.move(1)
	}
}

func BenchmarkMoveFull1000(b *testing.B) {
	benchmarkMoveFull(1000, b)
}

func BenchmarkMoveFull10000(b *testing.B) {
	benchmarkMoveFull(10000, b)
}

func benchmarkMoveAndSet(i int, b *testing.B) {
	r := newRing(i)
	for n := 0; n < b.N; n++ {
		r.value = n
		r.move(1)
	}
}

func BenchmarkMoveAndSet1000(b *testing.B) {
	benchmarkMoveAndSet(1000, b)
}

func BenchmarkMoveAndSet10000(b *testing.B) {
	benchmarkMoveAndSet(10000, b)
}
