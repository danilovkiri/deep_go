package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	data  []Task
	index map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{
		data:  []Task{},
		index: make(map[int]int),
	}
}

func (s *Scheduler) AddTask(t Task) {
	if _, exists := s.index[t.Identifier]; exists {
		return // duplicate identifiers not allowed
	}
	s.data = append(s.data, t)
	i := len(s.data) - 1
	s.index[t.Identifier] = i
	s.heapifyUp(i)
}

func (s *Scheduler) heapifyUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if s.data[i].Priority <= s.data[parent].Priority {
			break
		}
		s.swap(i, parent)
		i = parent
	}
}

func (s *Scheduler) heapifyDown(i int) {
	last := len(s.data) - 1
	for {
		left := 2*i + 1
		right := 2*i + 2
		largest := i

		if left <= last && s.data[left].Priority > s.data[largest].Priority {
			largest = left
		}
		if right <= last && s.data[right].Priority > s.data[largest].Priority {
			largest = right
		}
		if largest == i {
			break
		}
		s.swap(i, largest)
		i = largest
	}
}

func (s *Scheduler) swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
	s.index[s.data[i].Identifier] = i
	s.index[s.data[j].Identifier] = j
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	i, exists := s.index[taskID]
	if !exists {
		return
	}
	oldPriority := s.data[i].Priority
	s.data[i].Priority = newPriority
	if newPriority > oldPriority {
		s.heapifyUp(i)
	} else {
		s.heapifyDown(i)
	}
}

func (s *Scheduler) GetTask() Task {
	if len(s.data) == 0 {
		return Task{}
	}
	maxTask := s.data[0]
	last := len(s.data) - 1
	s.swap(0, last)
	s.data = s.data[:last]
	delete(s.index, maxTask.Identifier)
	if len(s.data) > 0 {
		s.heapifyDown(0)
	}
	return maxTask
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1.Identifier, task.Identifier)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
