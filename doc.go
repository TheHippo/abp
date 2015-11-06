// Package abp implements a buffer pool where the size
// of the buffers gets averages over all returned buffers.

// It uses fixed sized channel also known as leaky buffers
// (http://golang.org/doc/effective_go.html#leaky_buffer).

// The calculation of the average size uses a ring and will take
// place in a separate go routine.

package abp
