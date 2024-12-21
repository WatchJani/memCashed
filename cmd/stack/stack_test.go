package stack

import "testing"

const TEST_INPUT = 1000

func BenchmarkCustom(b *testing.B) {
	b.StopTimer()
	stack := New[int](TEST_INPUT)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for index := range TEST_INPUT {
			stack.Push(index)
		}

		stack.Clear()
	}
}
