package gqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
)

// An item is something we manage in a priority queue.
type item[T any] struct {
	value T

	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue[T any] struct {
	capacity    int
	compareFunc func(src T, dst T) bool
	queue       []item[T]
}

func (pq priorityQueue[T]) Len() int { return len(pq.queue) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq.compareFunc(pq.queue[i].value, pq.queue[j].value)
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
	pq.queue[i].index = i
	pq.queue[j].index = j
}

func (pq *priorityQueue[T]) Push(x any) {
	n := len(pq.queue)
	value := x.(T)
	item := item[T]{
		value: value,
		index: n,
	}
	pq.queue = append(pq.queue, item)
}

func (pq *priorityQueue[T]) Pop() any {
	old := pq.queue
	n := len(old)
	if n == 0 {
		return nil
	}
	i := old[n-1]
	i.index = -1
	pq.queue = old[0 : n-1]
	return i.value
}

type PriorityQueue[T any] struct {
	mu sync.Mutex
	pq priorityQueue[T]
}

var ErrEmptyQueue = errors.New("empty queue")

func NewPriorityQueue[T any](capacity int, compareFunc func(src T, dst T) bool) *PriorityQueue[T] {
	pq := priorityQueue[T]{
		capacity:    capacity,
		compareFunc: compareFunc,
		queue:       make([]item[T], 0, capacity),
	}
	return &PriorityQueue[T]{
		mu: sync.Mutex{},
		pq: pq,
	}
}

func (p *PriorityQueue[T]) Enqueue(elem T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pq.capacity > 0 {
		if p.pq.Len() < p.pq.capacity {
			// if not reached, just return
			heap.Push(&p.pq, elem)
			return
		}

		// if the capacity is reached, we have to take out the last element
		_ = heap.Remove(&p.pq, p.pq.capacity-1)
	}

	// Then push
	heap.Push(&p.pq, elem)
}

func (p *PriorityQueue[T]) Dequeue() (T, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.pq.Len() <= 0 {
		return *new(T), ErrEmptyQueue
	}

	ret, ok := heap.Pop(&p.pq).(T)
	if !ok {
		var t T
		// This should never happen
		return t, fmt.Errorf("value type %T popped out is not %T", ret, *new(T))
	}
	return ret, nil
}
