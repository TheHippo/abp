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

func (r *ring) move(n int) *ring {
	if r == nil {
		return nil
	}
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
