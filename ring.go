package abp

type ring struct {
	value      int
	prev, next *ring
}

func newRing(size int) *ring {
	if size < 1 {
		return nil
	}

	r := new(ring)
	p := r
	for i := 1; i < size; i++ {
		p.next = &ring{
			prev: p,
		}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

func (r *ring) Len() int {
	n := 0
	if r != nil {
		n++
		for p := r.next; p != r; p = p.next {
			n++
		}
	}
	return n
}

func (r *ring) Move(n int) *ring {
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}
