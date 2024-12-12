package server

import "testing"

// 0.2255 ns/op
func BenchmarkDecoder(b *testing.B) {
	b.StopTimer()
	payload := []byte{83, 11, 105, 203, 112, 126, 4, 0, 0, 0}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Decode(payload)
	}
}
