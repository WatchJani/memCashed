package stack

import (
	"errors"
)

// ErrEmpty represents the error when trying to pop from an empty stack.
var ErrEmpty = errors.New("stack is empty")

// Stack is a generic stack data structure that can store any type of value.
type Stack[T any] struct {
	store []T // A slice used to store the elements of the stack.
}

// New creates and returns a new stack with the specified capacity.
func New[T any](capacity int) Stack[T] {
	return Stack[T]{
		store: make([]T, 0, capacity), // Create an empty slice with the given capacity.
	}
}

// Push adds a new element to the top of the stack.
func (s *Stack[T]) Push(value T) {
	s.store = append(s.store, value) // Append the value to the slice, which represents the stack.
}

// Clear removes all elements from the stack, effectively resetting it.
func (s *Stack[T]) Clear() {
	s.store = s.store[:0] // Resize the slice to zero length, clearing all elements.
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns an error.
func (s *Stack[T]) Pop() (T, error) {
	if len(s.store) == 0 {
		var zeroValue T            // Return the zero value of type T.
		return zeroValue, ErrEmpty // Return the empty stack error.
	}
	// Get the top element (last element of the slice).
	value := s.store[len(s.store)-1]
	// Resize the slice to remove the top element.
	s.store = s.store[:len(s.store)-1]
	return value, nil // Return the popped value and nil error.
}

// IsEmpty checks if the stack is empty and returns a boolean result.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.store) == 0 // Return true if the stack is empty, false otherwise.
}

// Peek returns the top element of the stack without removing it.
// If the stack is empty, it returns an error.
func (s *Stack[T]) Peek() (T, error) {
	if len(s.store) == 0 {
		var zeroValue T            // Return the zero value of type T.
		return zeroValue, ErrEmpty // Return the empty stack error.
	}

	// Return the top element without removing it from the stack.
	return s.store[len(s.store)-1], nil
}
