package models

import (
	"errors"
	"slices"
)

/* DEFINITIONS */

var (
	ErrMaxCapacity = errors.New("Max capacity reached")
	ErrEmpty       = errors.New("Stack is empty")
)

type Stack[T any] struct {
	data []T
	max  uint
	last uint
}

/* FUNCTIONS */

func NewStack[T any](size uint) Stack[T] {
	return Stack[T]{
		data: make([]T, 0, size),
		max:  size,
		last: 0,
	}
}

func (s *Stack[T]) Push(v T) error {
	if s.last >= s.max-1 {
		return ErrMaxCapacity
	}

	s.data = append(s.data, v)
	s.last += 1
	return nil
}

func (s *Stack[T]) Pop() (T, error) {
	var empty T
	if s.last == 0 {
		return empty, ErrEmpty
	}

	val := s.data[s.last-1]
	s.data = slices.Delete(s.data, int(s.last-1), int(s.last))
	s.last -= 1

	return val, nil
}
