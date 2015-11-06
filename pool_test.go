package abp

import (
	"bytes"
	"runtime"
	"sync"

	"testing"
)

func TestPoolBasic(t *testing.T) {
	size := 10
	backlog := 20
	expected := 10
	overflow := 1.0

	content := []byte("foobar")
	p := NewAveragedBufferPool(size, backlog, expected, overflow)
	defer p.Close()
	for i := 0; i < backlog; i++ {
		buf := p.Get()
		buf.Write(content)
		_ = buf.String()
		p.Put(buf)
	}

	for i := 0; i < backlog; i++ {
		go func() {
			buf := p.Get()
			buf.Write(content)
			p.Put(buf)
		}()
	}
}

func TestPoolAverage(t *testing.T) {
	// this test is a little bit pointless, because it is impossible to predict how much
	// the go runtime will grow an buffer after writing over capacity
	size := 10
	backlog := 20
	expected := 20
	overflow := 1.0

	content := []byte("01234567890123456789")
	p := NewAveragedBufferPool(size, backlog, expected, overflow)
	defer p.Close()

	ca := p.currentAverage
	if ca != expected {
		t.Errorf("Expected current average to be %d but was %d", expected, ca)
	}

	cs := p.currentSum
	if cs != expected*backlog {
		t.Errorf("Expected current sum to be %d but was %s", expected*backlog, cs)
	}

	for i := 0; i < backlog; i++ {
		buf := p.Get()
		buf.Write(content)
		_ = buf.String()
		p.Put(buf)
	}

	// time.Sleep(1 * time.Second)
	p.lock.RLock()
	ca = p.currentAverage
	p.lock.RUnlock()
	if ca != len(content) {
		t.Errorf("Expected current average to be %d but was %d", len(content), ca)
	}

}

func TestCap(t *testing.T) {
	size := 10
	backlog := 20
	expected := 10
	overflow := 1.0

	runner := runtime.NumCPU()
	var wg sync.WaitGroup

	content := []byte("01234567890123456789")
	p := NewAveragedBufferPool(size, backlog, expected, overflow)
	defer p.Close()

	if cap(p.bufferChannel) != size {
		t.Errorf("Expected buffer size to be %d but was %d", size, cap(p.bufferChannel))
	}

	if len(p.bufferChannel) != 0 {
		t.Errorf("Expected buffer channel to be empty, but has %d buffers", len(p.bufferChannel))
	}

	wg.Add(runner)

	// just take buffers out

	for i := 0; i < runner; i++ {
		go func() {
			defer wg.Done()
			buffers := make([]*bytes.Buffer, backlog)
			for j := 0; j < backlog; j++ {
				buf := p.Get()
				buf.Write(content)
				_ = buf.String()
				buffers[j] = buf
				// p.Put(buf)
			}
		}()
	}
	wg.Wait()

	if len(p.bufferChannel) != 0 {
		t.Errorf("Channel should be empty, but there are %d elements", len(p.bufferChannel))
	}

	// take out and push all of them back in

	var wait sync.WaitGroup
	wg.Add(runner)
	wait.Add(runner)
	for i := 0; i < runner; i++ {
		go func() {
			buffers := make([]*bytes.Buffer, backlog)
			for j := 0; j < backlog; j++ {
				buf := p.Get()
				buf.Write(content)
				_ = buf.String()
				buffers[j] = buf
				// p.Put(buf)
			}
			wait.Done()
			wait.Wait()
			for _, buf := range buffers {
				p.Put(buf)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(p.bufferChannel) != size {
		t.Errorf("Channel should be full, but there just %d elements", len(p.bufferChannel))
	}

}
