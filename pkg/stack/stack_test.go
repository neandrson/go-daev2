package stack_test

import (
	"testing"

	"github.com/neandrson/go-daev2/pkg/stack"
)

func TestStack(t *testing.T) {
	stack := stack.NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	popped := stack.Pop()
	if popped != 2 {
		t.Errorf("Expected popped element to be 2, but got %d", popped)
	}
}
