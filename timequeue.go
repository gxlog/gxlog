package gxlog

import (
	"time"
)

type timeQueue struct {
	duration time.Duration
	slice    []time.Time
	cap      int
	begin    int
	end      int
}

func newTimeQueue(duration time.Duration, size int) *timeQueue {
	return &timeQueue{
		duration: duration,
		slice:    make([]time.Time, size+1),
		cap:      size + 1,
	}
}

func (this *timeQueue) Enqueue(clock time.Time) bool {
	this.dequeueExpired(clock)
	if this.full() {
		return false
	}
	this.slice[this.end] = clock
	this.end++
	this.end %= this.cap
	return true
}

func (this *timeQueue) dequeueExpired(clock time.Time) {
	if !this.empty() && clock.Sub(this.slice[this.begin]) >= this.duration {
		this.begin++
		this.begin %= this.cap
	}
}

func (this *timeQueue) length() int {
	return (this.end + this.cap - this.begin) % this.cap
}

func (this *timeQueue) empty() bool {
	return this.begin == this.end
}

func (this *timeQueue) full() bool {
	return this.length() == this.cap-1
}
