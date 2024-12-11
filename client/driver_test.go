package client

import "testing"

// 212ns
func BenchmarkDriver(b *testing.B) {
	b.StopTimer()

	key := []byte("janko")
	value := []byte("super dan za pobjedu\r\n")
	ttl := 3 * 60 * 60

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Req(key, value, ttl)
	}
}
