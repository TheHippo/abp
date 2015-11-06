package abp

import (
	"bytes"
	"sync/atomic"
)

type AveragedBufferPool struct {
	currentAverage *int64
	max            *int64
	backlog        int
	currentSum     int
	overflow       float32
	ring           *ring
	bufferChannel  chan *bytes.Buffer
	calcChannel    chan int
	closeCalc      chan struct{}
}

type int64Wrapper struct {
	x int64
}

func NewAveragedBufferPool(size, backlog, expectedSize int, overflow float32) *AveragedBufferPool {
	r := newRing(backlog)
	for i := 0; i < backlog; i++ {
		r.value = expectedSize
		r = r.move(1)
	}
	abp := &AveragedBufferPool{
		backlog:    backlog,
		currentSum: backlog * expectedSize,
		currentAverage: &(&int64Wrapper{
			x: int64(expectedSize),
		}).x,
		max: &(&int64Wrapper{
			x: int64(float32(expectedSize) * overflow),
		}).x,
		overflow:      overflow,
		ring:          r,
		bufferChannel: make(chan *bytes.Buffer, size),
		calcChannel:   make(chan int),
		closeCalc:     make(chan struct{}),
	}
	go abp.calc()
	return abp
}

func (abp *AveragedBufferPool) Close() {
	close(abp.closeCalc)
}

func (abp *AveragedBufferPool) calc() {
	for {
		select {
		case <-abp.closeCalc:
			return
		case lastSize := <-abp.calcChannel:
			// update current sum, write new value to ring
			abp.currentSum -= abp.ring.value
			abp.currentSum += lastSize
			abp.ring.value = lastSize
			abp.ring.move(1)

			// calculate average and max
			ca := int64(abp.currentSum / abp.backlog)
			atomic.StoreInt64(abp.currentAverage, ca)
			atomic.StoreInt64(abp.max, int64(float32(ca)+abp.overflow))
		}
	}
}

func (abp *AveragedBufferPool) Get() *bytes.Buffer {
	select {
	case b := <-abp.bufferChannel:
		return b
	default:
		return bytes.NewBuffer(make([]byte, 0, atomic.LoadInt64(abp.currentAverage)))
	}
}

func (abp *AveragedBufferPool) Put(buf *bytes.Buffer) {
	// rewind buffer to start
	buf.Reset()

	capacity := cap(buf.Bytes())

	if int64(capacity) > atomic.LoadInt64(abp.max) {
		// buffer was larger the average. create new one
		buf = bytes.NewBuffer(make([]byte, 0, atomic.LoadInt64(abp.currentAverage)))
	}

	select {
	case abp.bufferChannel <- buf:
		// return buffer into channel
	default:
		// channel is full, throw it away
	}

	// send value to calculate new averages
	abp.calcChannel <- capacity
}
