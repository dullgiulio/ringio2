package main

import "fmt"

type ring struct{
	mux sync.Mutex
	cond sync.Cond
	data [][]byte
	pos int
	rdr int // position of slowest reader
}

func newRing(size int) *ring {
	r := &ring{
		data: make([][]byte, size, size)
	}
	r.cond.L = &r.mux
	return r
}

func (r *ring) put(data []byte) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.pos >= r.rdr {
		r.cond.Broadcast()
	}
	r.data[r.pos] = data
	r.pos = (r.pos + 1) % len(r.data)
}

type ringReader struct{
	ring *ring
	pos int
}

func newRingReader(r *ring) *ringReader {
	return &ringReader{
		ring: r,
	}
}

func (r *ringReader) get() []byte {
	r.ring.mux.Lock()
	defer r.ring.mux.Unlock()

	r.checkLast()
	for dist(r.pos, r.ring.pos, r.ring.rdr, len(r.ring.data)) <= 0 {
		r.cond.Wait()
	}
	data := r.ring.data[r.pos]
	r.pos = (r.pos + 1) % len(r.ring.data)
	return data
}

func (r *ringReader) checkLast() {
	if r.ring.pos < r.ring.rdr && r.pos > r.ring.pos && r.pos < r.ring.rdr {
		r.ring.rdr = r.pos
		return
	}
}

func dist(a, b, zero, size int) int {
	if (zero <= a && zero <= b) || (zero > a && zero > b) {
		return b - a
	}
	if a <= zero && b >= zero {
		return size - b + a
	}
	panic(fmt.Sprintf("dist: a = %d, b = %d, zero = %d, size = %d", a, b, zero, size))
}

func main() {
}
