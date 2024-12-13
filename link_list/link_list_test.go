package link_list

import "testing"

// 33.26 ns/op
func BenchmarkInset(b *testing.B) {
	b.StopTimer()

	f := new(DLL)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Inset(Value{})
	}
}
