package main

import (
	"testing"
)

func TestFibonacciIterative(t *testing.T) {
	fibData := []int{
		0, 1, 1, 2, 3, 5, 8, 13, 21,
	}
	for i := 0; i < len(fibData); i++ {
		if FibonacciIterative(i) != fibData[i] {
			t.Errorf("FibonacciIterative(%d) = %d", i, FibonacciIterative(i))
		}
	}
}
