package abp

type ring struct {
	value      int
	prev, next *ring
}

func newRing(size int) *ring {
	if size < 1 {
		return nil
	}

	r := &ring{}
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

func (r *ring) len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.next; p != r; p = p.next {
			n++
		}
	}
	return n
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

// func (r *ring) String() string {
// 	buf := bytes.NewBuffer(nil)
// 	n := r.len()
// 	buf.WriteString(fmt.Sprintf("[%d] ", n))
// 	p := r
// 	for ; n > 0; n-- {
// 		buf.WriteString(fmt.Sprintf("%d, ", p.value))
// 		p = p.move(1)
// 	}
// 	return buf.String()
// }
