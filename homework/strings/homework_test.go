package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"testing"
	"unsafe"
)

type COWBuffer struct {
	data []byte
	refs *int
}

func NewCOWBuffer(data []byte) COWBuffer {
	b := COWBuffer{
		data: data,     // not a copy of `data` here since tests require the original and buffer's arrays to be the same in terms of memory
		refs: new(int), // is a pointer since copies of buffer must share this value, *refs > 0 if active copies are present
	}

	runtime.SetFinalizer(&b, func(b *COWBuffer) {
		b.Close()
	})

	return b
}

func (b *COWBuffer) Clone() COWBuffer {
	*b.refs++

	return COWBuffer{
		data: b.data,
		refs: b.refs,
	}
}

func (b *COWBuffer) Close() {
	*b.refs--
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b == nil || index < 0 || index >= len(b.data) {
		return false
	}

	// if copies were made, then we make a new data and new refs before changing the value
	if *b.refs > 0 {
		*b.refs--
		newData := make([]byte, len(b.data))
		copy(newData, b.data)
		nb := NewCOWBuffer(newData)
		*b = nb
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
	fmt.Println("buffer.data", buffer.data)
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	fmt.Println("copy1.data", copy1.data)
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))
	fmt.Println("copy2.data", copy2.data)

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
