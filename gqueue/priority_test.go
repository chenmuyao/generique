package gqueue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue(3, func(src, dst int) bool {
		return src < dst
	})

	pq.Enqueue(3)
	pq.Enqueue(5)
	ret, err := pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 3, ret)

	pq.Enqueue(1)
	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, ret)

	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 5, ret)
	_, err = pq.Dequeue()
	assert.ErrorIs(t, err, ErrEmptyQueue)

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := range 3 {
		go func(idx int) {
			defer wg.Done()
			pq.Enqueue(idx)
		}(i)
	}
	wg.Wait()

	pq.Enqueue(3)
	assert.Equal(t, 3, len(pq.pq.queue))
	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, ret)

	// readd 1
	pq.Enqueue(1)
	assert.Equal(t, 3, len(pq.pq.queue))
	// adding 0 should be ignored
	pq.Enqueue(0)

	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, ret)
	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 2, ret)
	ret, err = pq.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 3, ret)
}
