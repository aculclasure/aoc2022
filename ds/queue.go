package ds

import "sync"

type Queue[T any] struct {
	mtx  *sync.Mutex
	vals []T
}

func (q *Queue[T]) Enqueue(val T) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.vals = append(q.vals, val)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if len(q.vals) == 0 {
		var zero T
		return zero, false
	}
	front := q.vals[0]
	q.vals = q.vals[1:]
	return front, true
}

func (q *Queue[T]) Size() int {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	return len(q.vals)
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		mtx: new(sync.Mutex),
	}
}
