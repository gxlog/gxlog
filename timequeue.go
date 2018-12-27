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

func (queue *timeQueue) Enqueue(clock time.Time) bool {
	queue.dequeueExpired(clock)
	if queue.full() {
		return false
	}
	queue.slice[queue.end] = clock
	queue.end++
	queue.end %= queue.cap
	return true
}

func (queue *timeQueue) dequeueExpired(clock time.Time) {
	if !queue.empty() && clock.Sub(queue.slice[queue.begin]) >= queue.duration {
		queue.begin++
		queue.begin %= queue.cap
	}
}

func (queue *timeQueue) length() int {
	return (queue.end + queue.cap - queue.begin) % queue.cap
}

func (queue *timeQueue) empty() bool {
	return queue.begin == queue.end
}

func (queue *timeQueue) full() bool {
	return queue.length() == queue.cap-1
}
