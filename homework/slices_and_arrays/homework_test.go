package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue[T int | int8 | int16 | int32 | int64] struct {
	values            []T
	size, front, rear int
}

func NewCircularQueue[T int | int8 | int16 | int32 | int64](size int) CircularQueue[T] {
	return CircularQueue[T]{
		values: make([]T, size),
		size:   size,
		front:  -1,
		rear:   -1,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	if q.Empty() {
		q.front = 0
	}

	q.rear = (q.rear + 1) % q.size
	q.values[q.rear] = value
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	if q.front == q.rear {
		q.front = -1
		q.rear = -1
		return true
	}

	q.front = (q.front + 1) % q.size
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	if q.Empty() {
		return -1
	}
	res := q.values[q.front]
	return res

	//if q.front == q.rear {
	//	result := q.values[q.front]
	//	q.front = -1
	//	q.rear = -1
	//	return result
	//}
	//
	//result := q.values[q.front]
	//q.front = (q.front + 1) % q.size
	//return result
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	if q.Empty() {
		return -1
	}
	res := q.values[q.rear]
	return res

	//if q.front == q.rear {
	//	result := q.values[q.front]
	//	q.front = -1
	//	q.rear = -1
	//	return result
	//}
	//
	//result := q.values[q.rear]
	//if q.rear == 0 {
	//	q.rear = 5
	//} else {
	//	q.rear -= 1
	//}
	//return result
}

func (q *CircularQueue[T]) Empty() bool {
	return q.rear == -1 && q.front == -1
}

func (q *CircularQueue[T]) Full() bool {
	return (q.front == 0 && q.rear == q.size-1) || q.front == q.rear+1
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
