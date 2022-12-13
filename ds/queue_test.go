package ds_test

import (
	"testing"

	"github.com/aculclasure/aoc2022/ds"
)

func TestQueue_DequeueFromEmptyQueueReturnsFalse(t *testing.T) {
	t.Parallel()
	q := ds.NewQueue[int]()
	_, ok := q.Dequeue()
	if ok {
		t.Error("want false, got true")
	}
}

func TestQueue_DequeueReturnsFrontOfQueue(t *testing.T) {
	t.Parallel()
	q := ds.NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	want := 1
	got, ok := q.Dequeue()
	if !ok {
		t.Fatal("want true status, got false indicating queue is empty")
	}
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestQueue_DequeueDecreasesQueueSize(t *testing.T) {
	t.Parallel()
	q := ds.NewQueue[int]()
	q.Enqueue(1)
	_, ok := q.Dequeue()
	if !ok {
		t.Fatal("want true status, got false indicating queue is empty")
	}
	want := 0
	got := q.Size()
	if want != got {
		t.Errorf("want size %d, got size %d", want, got)
	}
}

func TestQueue_EnqueueIncreasesQueueSize(t *testing.T) {
	t.Parallel()
	q := ds.NewQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	want := 3
	got := q.Size()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
