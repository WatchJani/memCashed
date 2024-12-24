package memory_allocator

import (
	"testing"
	"unsafe"
)

func BenchmarkAllocator(b *testing.B) {
	b.StopTimer()

	allocator := New(5 * 1024 * 1024)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		allocator.AllocateBlock()
	}
}

func BenchmarkUnsafePointer(b *testing.B) {
	b.StopTimer()
	data := make([]byte, 10)
	ptr := unsafe.Pointer(&data[0])

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = unsafe.Slice((*byte)(ptr), 64)
	}
}
