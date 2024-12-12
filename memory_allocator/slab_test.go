package memory_allocator

import (
	"log"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	binarySearch := func(dataSize int) int {
		slabs := []int{
			64,
			128,
			256,
			512,
			1024,
			2048,
			4096,
			8192,
			16384,
			32768,
			65536,
			131072,
			262144,
			524288,
			1048576,
		}

		low, high := 0, len(slabs)-1
		result := high

		for low <= high {
			mid := low + (high-low)/2
			if slabs[mid] >= dataSize {
				result = mid
				high = mid - 1
			} else {
				low = mid + 1
			}
		}

		return result
	}

	testCase := []struct {
		input   int
		results int
	}{
		{1048596, 14},
		{10, 0},
		{128, 1},
	}

	for index, test := range testCase {
		if output := binarySearch(test.input); output != test.results {
			t.Errorf("%d | get value %d | expected value %d", index, output, test.results)
		}
	}
}

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
