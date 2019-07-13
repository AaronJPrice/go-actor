package actor

import (
	"testing"
)

//==============================================================================
// Benchmarks
//==============================================================================
func BenchmarkOneActorManyRequests(b *testing.B) {
	actorReg := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		actorReg.Execute(1, fun)
	}
}

func BenchmarkManyActorsOneRequest(b *testing.B) {
	actorReg := New()
	var i int64

	b.ResetTimer()
	for i = 0; i < int64(b.N); i++ {
		actorReg.Execute(i, fun)
	}
}

func BenchmarkManyActorsManyRequest(b *testing.B) {
	actorReg := New()
	var o, j int64
	o = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j = 0; j < o; j++ {
			actorReg.Execute(j, fun)
		}
	}
}

//==============================================================================
// Utilities
//==============================================================================
func fun() interface{} {
	return 1 + 1
}
