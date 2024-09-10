package main

import "testing"

func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = hodgepodge()
	}
}
