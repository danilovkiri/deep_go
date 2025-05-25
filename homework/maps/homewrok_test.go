package main

import (
	"cmp"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

const (
	less    = -1
	equal   = 0
	greater = 1
)

type Node[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *Node[K, V]
	right *Node[K, V]
}

type OrderedMap[K cmp.Ordered, V any] struct {
	root *Node[K, V]
	size int
	cmp  func(x, y K) int
}

func NewOrderedMap[K cmp.Ordered, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		root: nil,
		size: 0,
		cmp:  cmp.Compare[K],
	}
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	if node := findNode(m.root, key, m.cmp); node != nil {
		return node.value, true
	}

	var zero V
	return zero, false
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.root = m.insertRecursive(m.root, key, value)
}

func (m *OrderedMap[K, V]) insertRecursive(root *Node[K, V], key K, value V) *Node[K, V] {
	if root == nil {
		root = &Node[K, V]{key: key, value: value}
		m.size++
		return root
	}

	switch m.cmp(key, root.key) {
	case less:
		root.left = m.insertRecursive(root.left, key, value)
	case greater:
		root.right = m.insertRecursive(root.right, key, value)
	case equal:
		root.value = value
	}

	return root
}

func (m *OrderedMap[K, V]) Erase(key K) {
	m.root = m.remove(m.root, key)
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	if node := findNode(m.root, key, m.cmp); node != nil {
		return true
	}
	return false
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func (m *OrderedMap[K, V]) ForEach(action func(kye K, value V)) {
	traverse(m.root, action)
}

func (m *OrderedMap[K, V]) remove(root *Node[K, V], key K) *Node[K, V] {
	if root == nil {
		return nil
	}

	switch m.cmp(key, root.key) {
	case less:
		root.left = m.remove(root.left, key)
	case greater:
		root.right = m.remove(root.right, key)
	case equal:
		if root.left == nil {
			m.size--
			return root.right
		} else if root.right == nil {
			m.size--
			return root.left
		}

		// find the min node of the right subtree
		minNode := root.right
		for minNode != nil && minNode.left != nil {
			minNode = minNode.left
		}

		root.key = minNode.key
		root.value = minNode.value
		root.right = m.remove(root.right, minNode.key)
	}
	return root
}

func findNode[K cmp.Ordered, V any](root *Node[K, V], key K, cmp func(x, y K) int) *Node[K, V] {
	if root == nil {
		return nil
	}

	switch cmp(key, root.key) {
	case less:
		return findNode[K, V](root.left, key, cmp)
	case greater:
		return findNode[K, V](root.right, key, cmp)
	default:
		return root
	}
}

func traverse[K cmp.Ordered, V any](node *Node[K, V], action func(key K, value V)) {
	if node == nil {
		return
	}
	traverse(node.left, action)
	action(node.key, node.value)
	traverse(node.right, action)
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
