package abp

import (
	"bytes"
	"sync"
)

type AveragedBufferPool struct {
	currentAverage int
	max            int
	backlog        int
	currentSum     int
	overflow       float64
	ring           *ring
	bufferChannel  chan *bytes.Buffer
	calcChannel    chan int
	closeCalc      chan struct{}
	lock           sync.RWMutex
}

func NewAveragedBufferPool(size, backlog, expectedSize int, overflow float64) *AveragedBufferPool {
	r := newRing(backlog)
	for i := 0; i < backlog; i++ {
		r.value = expectedSize
		r = r.move(1)
	}
	abp := &AveragedBufferPool{
		backlog:        backlog,
		currentSum:     backlog * expectedSize,
		currentAverage: expectedSize,
		max:            int(float64(expectedSize) * overflow),
		overflow:       overflow,
		ring:           r,
		bufferChannel:  make(chan *bytes.Buffer, size),
		calcChannel:    make(chan int),
		closeCalc:      make(chan struct{}),
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
			abp.lock.Lock()
			abp.currentSum -= abp.ring.value
			abp.currentSum += lastSize
			abp.ring.value = lastSize
			abp.ring = abp.ring.move(1)
			abp.currentAverage = abp.currentSum / abp.backlog
			abp.max = int(float64(abp.currentAverage) * abp.overflow)
			abp.lock.Unlock()
		}
	}
}

func (abp *AveragedBufferPool) Get() *bytes.Buffer {
	select {
	case b := <-abp.bufferChannel:
		return b
	default:
		abp.lock.RLock()
		buf := bytes.NewBuffer(make([]byte, 0, abp.currentAverage))
		abp.lock.RUnlock()
		return buf
	}
}

func (abp *AveragedBufferPool) Put(buf *bytes.Buffer) {
	// rewind buffer to start
	buf.Reset()

	capacity := cap(buf.Bytes())

	abp.lock.RLock()

	if capacity > abp.max {
		// buffer was larger the average. create new one
		buf = bytes.NewBuffer(make([]byte, 0, abp.currentAverage))
	}

	abp.lock.RUnlock()

	select {
	case abp.bufferChannel <- buf:
		// return buffer into channel
	default:
		// channel is full, throw it away
	}

	// send value to calculate new averages
	abp.calcChannel <- capacity
}
