package memory_allocator

import "testing"

func BenchmarkAllocator(b *testing.B) {
	b.StopTimer()

	allocator := New(5 * 1024 * 1024)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		allocator.AllocateBlock()
	}
}
