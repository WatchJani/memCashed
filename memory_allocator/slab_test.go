package memory_allocator

import (
	"log"
	"testing"
)

func BenchmarkSlabAllocator(b *testing.B) {
	b.StopTimer()

	allocator := New(5 * 1024 * 1024) //5 MiB
	s := NewSlab(64, allocator)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.AllocateMemory()
		if err != nil {
			log.Println(err)
		}
	}
}
