package auth_service

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type RateLimiter struct {
	queue *list.List
	mutex *sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		queue: list.New(),
		mutex: &sync.Mutex{},
	}
}

const (
	OneMinute = 60000000
)

var (
	ipMap          = make(map[string]*RateLimiter)
	blockedIPs     = make(map[string]bool)
	blockedIPMutex sync.Mutex
	ipMapMutex     sync.Mutex
)

func (r *RateLimiter) Enqueue(value time.Time) {
	r.mutex.Lock()
	r.queue.PushBack(value)
	r.mutex.Unlock()
}

func (r *RateLimiter) Dequeue() error {
	if r.Empty() {
		return errors.New("Queue empty")
	}

	r.mutex.Lock()
	ele := r.queue.Front()
	r.queue.Remove(ele)
	r.mutex.Unlock()
	return nil
}

func (r *RateLimiter) Front() (time.Time, error) {
	if r.Empty() {
		return time.Time{}, errors.New("queue empty")
	}

	r.mutex.Lock()
	val, ok := r.queue.Front().Value.(time.Time)
	r.mutex.Unlock()

	if !ok {
		return time.Time{}, errors.New("incorrect datatype")
	}

	return val, nil

}

func (r *RateLimiter) Size() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.queue.Len()
}

func (r *RateLimiter) Empty() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return (r.queue.Len() == 0)
}

func (c *RateLimiter) Format(currentTime time.Time) {
	for {
		if c.Empty() {
			return
		}

		if c.IsThresholdReached() {
			return
		}

		previousTime, err := c.Front()
		if err != nil {
			return
		}

		if time.Since(previousTime).Microseconds() > (int64(OneMinute)) {
			err = c.Dequeue()
			if err != nil {
				return
			}
		} else {
			break
		}
	}
}

func (r *RateLimiter) IsThresholdReached() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return (r.queue.Len()) > 50
}
