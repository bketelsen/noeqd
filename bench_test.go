package noeqd

import "testing"

func BenchmarkIdGeneration(b *testing.B) {
	for n := 0; n < b.N; n++ {
		nextId()
	}
}
