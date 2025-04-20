package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

type COWBuffer struct {
	data []byte
	refs *int
	mu   *sync.Mutex
}

func NewCOWBuffer(data []byte) COWBuffer {
	if len(data) == 0 {
		panic("data is empty") // data can be empty only in a buffer instance which was closed
	}

	b := COWBuffer{
		data: data,     // not a copy of `data` here since tests require the original and buffer's arrays to be the same in terms of memory
		refs: new(int), // is a pointer since copies of buffer must share this value, *refs > 0 if active copies are present
		mu:   &sync.Mutex{},
	}

	runtime.SetFinalizer(&b, func(b *COWBuffer) {
		b.Close()
	})

	return b
}

func (b *COWBuffer) Clone() COWBuffer {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.data == nil {
		panic("cannot copy closed buffer")
	}

	*b.refs++

	return COWBuffer{
		data: b.data,
		refs: b.refs,
		mu:   b.mu,
	}
}

func (b *COWBuffer) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.data == nil {
		return // if data == nil, then this copy is already closed, multiple Close() calls must not throw an error
	}

	*b.refs--
	// so that no one ever tries to use it
	b.data = nil
	b.refs = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.data == nil {
		panic("cannot update closed buffer")
	}

	if b == nil || index < 0 || index >= len(b.data) {
		return false
	}

	// if copies were made, then we make a new data, new refs and new mutex prior to changing the value
	if *b.refs > 0 {
		newData := make([]byte, len(b.data))
		copy(newData, b.data)
		*b.refs--

		b.data = newData
		b.refs = new(int)
		b.mu = &sync.Mutex{}
	}

	b.data[index] = value
	return true
}

func (b *COWBuffer) String() string {
	if len(b.data) == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}

func TestCOWBuffer2(t *testing.T) {
	var data []byte
	assert.Panics(t, func() {
		_ = NewCOWBuffer(data)
	})
}

func TestCOWBuffer3(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	buffer.Close()

	assert.Panics(t, func() {
		_ = buffer.Clone()
	})

	assert.Panics(t, func() {
		_ = buffer.Update(0, 'p')
	})
}
